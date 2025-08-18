#!/usr/bin/env python3
"""
export_property_json.py

Connects to a Postgres database, queries the `properties` table,
and outputs JSON objects in the target structure.

Now with .env support, verbose diagnostics, and DSN printing.
"""

import argparse
import json
import os
import sys
from typing import Any, Dict, List, Optional

try:
    import psycopg
    import psycopg.rows
except ImportError:
    print("ERROR: This script requires psycopg (v3). Install with: pip install psycopg[binary]")
    sys.exit(1)

try:
    from dotenv import load_dotenv
except ImportError:
    print("ERROR: This script requires python-dotenv. Install with: pip install python-dotenv")
    sys.exit(1)


def infer_doc_type(url: str) -> str:
    u = url.lower()
    if "title" in u or "deed" in u or "certificate" in u:
        return "title_deed"
    if "survey" in u or "report" in u or "map" in u:
        return "survey_report"
    if u.endswith((".jpg", ".jpeg", ".png", ".gif", ".webp")):
        return "photo"
    return "document"

def load_trello_mapping(conn: "psycopg.Connection", trello_urls: List[str]) -> Dict[str, str]:
    """
    Given a list of trello_urls, fetch storage_url mapping from the attachments table.
    Returns dict: {trello_url: storage_url}
    """
    if not trello_urls:
        return {}
    # Deduplicate
    uniq = sorted(set(trello_urls))
    placeholders = ", ".join(["%s"] * len(uniq))
    sql = f"SELECT trello_url, storage_url FROM attachments WHERE trello_url IN ({placeholders})"
    mapping: Dict[str, str] = {}
    with conn.cursor() as cur:
        cur.execute(sql, uniq)
        for trello_url, storage_url in cur.fetchall():
            mapping[trello_url] = storage_url
    return mapping

def row_to_json(row: Dict[str, Any], city: str, country: str, square_meter_price: float, url_map: Optional[Dict[str, str]] = None) -> Dict[str, Any]:
    location_en = row["location_en"]
    lot = row["lot"]
    owner = row["owner"]
    attachments: Optional[List[str]] = row.get("attachments")
    docs = []
    if attachments:
        for url in attachments:
            if url and isinstance(url, str):
                mapped = url_map.get(url) if url_map else None
                final_url = mapped if mapped else url
                docs.append({"type": infer_doc_type(final_url), "url": final_url})
    return {
        "name": f"{location_en}-{lot}",
        "lot": lot,
        "description": f"Historic family land in the {location_en} district, originally acquired by {owner} in the late 19th century.",
        "location": location_en,
        "city": city,
        "country": country,
        "area": float(row["area"]) if row.get("area") is not None else None,
        "area_unit": "square meters",
        "square_meter_price": float(square_meter_price),
        "shares": float(row["shares"]) if row.get("shares") is not None else None,
        "owner": owner,
        "possessed": (row["type"] == "Possessed"),
        "unsold": (row["status"] == "Unsold"),
        "organized": bool(row["is_organized"]),
        "effects": bool(row["is_effects"]),
        "documents": docs,
    }


def build_filters(args: argparse.Namespace) -> (str, List[Any]):
    clauses = [
        "type IN ('Normal','Possessed')",
        "status = 'Unsold'",
        "owner NOT IN ('Uknown', '')"
    ]
    params: List[Any] = []
    if args.id is not None:
        clauses.append("id = %s")
        params.append(args.id)
    if args.location is not None:
        clauses.append("location_en = %s")
        params.append(args.location)
    if args.lot is not None:
        clauses.append("lot = %s")
        params.append(args.lot)
    where_sql = " WHERE " + " AND ".join(clauses) if clauses else ""
    return where_sql, params


def fetch_rows(conn: "psycopg.Connection", args: argparse.Namespace) -> List[Dict[str, Any]]:
    where_sql, params = build_filters(args)
    if not (args.id is not None or (args.location and args.lot) or args.all):
        print("ERROR: specify one of --id, --location and --lot together, or --all", file=sys.stderr)
        sys.exit(2)
    sql = f"""
        SELECT
            id, board_id, card_id, name, location_ar, location_en,
            lot, type, status, owner, area, shares,
            is_organized, is_effects, labels, attachments, comments, updated_at
        FROM properties
        {where_sql}
        ORDER BY updated_at DESC, id DESC
    """
    with conn.cursor(row_factory=psycopg.rows.dict_row) as cur:
        cur.execute(sql, params)
        return list(cur.fetchall())



def is_url_like(val: Optional[str]) -> bool:
    if not val:
        return False
    return val.startswith("postgres://") or val.startswith("postgresql://")


def maybe_append_sslmode(dsn: str) -> str:
    # If DSN already specifies sslmode, leave it. Otherwise, for typical cloud hosts add require.
    if "sslmode=" in dsn or "sslmode=%" in dsn:
        return dsn
    # Heuristic: if host looks like a domain (not local) and not 127.0.0.1, add require
    # Note: psycopg key-value DSN: look for host=... token
    host_match = None
    for token in dsn.split():
        if token.lower().startswith("host="):
            host_match = token.split("=", 1)[1]
            break
    # For URL DSN, simple check
    if dsn.startswith(("postgres://", "postgresql://")):
        # crude: if contains localhost or 127.0.0.1, skip
        if "localhost" in dsn or "127.0.0.1" in dsn:
            return dsn
        return f"{dsn}?sslmode=require" if "?" not in dsn else (dsn + "&sslmode=require")
    else:
        if host_match and host_match not in ("localhost", "127.0.0.1"):
            return dsn + " sslmode=require"
    return dsn


def redact_pw(dsn_val: str) -> str:
    if not dsn_val:
        return ""
    out = []
    for token in dsn_val.split():
        if token.lower().startswith("password="):
            out.append("password=***")
        else:
            out.append(token)
    return " ".join(out)


def main():
    parser = argparse.ArgumentParser(description="Export property JSON from Postgres.")
    parser.add_argument("--dsn", help="Full psycopg DSN. If omitted, uses PG* env vars.", default=None)
    parser.add_argument("--id", type=int, help="Property id to export.")
    parser.add_argument("--location", help="Filter by location_en (use with --lot).")
    parser.add_argument("--lot", help="Filter by lot (use with --location).")
    parser.add_argument("--all", action="store_true", help="Export all rows matching default filter.")
    parser.add_argument("--city", default="Damascus", help="City to populate in JSON (default: Damascus).")
    parser.add_argument("--country", default="Syria", help="Country to populate in JSON (default: Syria).")
    parser.add_argument("--square-meter-price", type=float, default=10.0, help="Square meter price to include.")
    parser.add_argument("--out", help="Write output to a file. If omitted, prints to stdout.")
    parser.add_argument("--ndjson", action="store_true", help="Emit one JSON object per line instead of an array.")
    parser.add_argument("--verbose", "-v", action="store_true", help="Verbose diagnostics.")
    parser.add_argument("--print-dsn", action="store_true", help="Print constructed DSN and exit (password redacted).")
    args = parser.parse_args()

    # Load .env
    dotenv_path = os.path.join(os.path.dirname(__file__), ".env")
    if os.path.exists(dotenv_path):
        load_dotenv(dotenv_path=dotenv_path)

    dsn = args.dsn
    if not dsn:
        # 1) DATABASE_URL takes precedence if present
        db_url = os.getenv("DATABASE_URL")
        if db_url and is_url_like(db_url):
            dsn = db_url
        else:
            # 2) PGHOST may already be a full URL; if so, use it
            pghost = os.getenv("PGHOST")
            if is_url_like(pghost):
                dsn = pghost
            else:
                # 3) Build key=value DSN from PG* vars
                parts = []
                for env_key, dsn_key in [
                    ("PGHOST", "host"),
                    ("PGPORT", "port"),
                    ("PGDATABASE", "dbname"),
                    ("PGUSER", "user"),
                    ("PGPASSWORD", "password"),
                    ("PGSSLMODE", "sslmode"),
                ]:
                    val = os.getenv(env_key)
                    if val:
                        parts.append(f"{dsn_key}={val}")
                dsn = " ".join(parts) if parts else None
    if not dsn:
        # Also check RAILWAY_DATABASE_URL (common on Railway)
        rw_url = os.getenv("RAILWAY_DATABASE_URL")
        if rw_url and is_url_like(rw_url):
            dsn = rw_url
    if dsn:
        dsn = maybe_append_sslmode(dsn)


    if args.verbose or args.print_dsn:
        env_snapshot = {
            "PGHOST": os.getenv("PGHOST"),
            "PGPORT": os.getenv("PGPORT"),
            "PGDATABASE": os.getenv("PGDATABASE"),
            "PGUSER": os.getenv("PGUSER"),
            "PGPASSWORD": "***" if os.getenv("PGPASSWORD") else None,
            "PGSSLMODE": os.getenv("PGSSLMODE"),
        }
        print("ENV snapshot:", env_snapshot, file=sys.stderr)
        print("Effective DSN:", redact_pw(dsn or ""), file=sys.stderr)
        if args.print_dsn:
            sys.exit(0)

    if not dsn:
        print("ERROR: Provide a --dsn or set PGHOST/PGPORT/PGDATABASE/PGUSER/PGPASSWORD.", file=sys.stderr)
        sys.exit(2)

    try:
        with psycopg.connect(dsn, connect_timeout=10) as conn:
            rows = fetch_rows(conn, args)

            # --- NEW: gather Trello URLs and load mappings ---
            trello_urls: List[str] = []
            for r in rows:
                atts = r.get("attachments")
                if atts and isinstance(atts, list):
                    for u in atts:
                        if isinstance(u, str) and u.startswith("https://trello.com"):
                            trello_urls.append(u)

            url_map = load_trello_mapping(conn, trello_urls)
    except Exception as e:
        print("DB ERROR:", e, file=sys.stderr)
        sys.exit(1)

    # Now build JSON using the url_map
    objs = [
        row_to_json(
            r,
            city=args.city,
            country=args.country,
            square_meter_price=args.square_meter_price,
            url_map=url_map,
        )
        for r in rows
    ]

    if args.ndjson:
        text = "\n".join(json.dumps(o, ensure_ascii=False) for o in objs)
    else:
        text = json.dumps(objs, ensure_ascii=False, indent=2)

    if args.out:
        with open(args.out, "w", encoding="utf-8") as f:
            f.write(text)
        print(f"Wrote {len(objs)} record(s) to {args.out}")
    else:
        print(text)


if __name__ == "__main__":
    main()
