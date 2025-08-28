## Prerequisites

- Python 3.11+
- Dependencies in `requirements.txt`

---

## Installation

1. **Install dependencies:**

```bash
python3 -m venv venv
source venv/bin/activate 
pip3 install -r requirements.txt
```

2. **Set up environment variables:**

- Copy `.env.example` to `.env`
- Edit `.env` with your API keys and preferences.

---

## Run Locally

### Streamlit Web Application

```bash
streamlit run app.py
```

### PDF Report Generation

Generate a comprehensive inheritance report in PDF format:

```bash
python3 run_report.py [output_filename.pdf]
```

**Examples:**
```bash
# Generate report with default name
python3 run_report.py

# Generate report with custom name
python3 run_report.py family_inheritance_report_2025.pdf
```

**Report Contents:**
- Family demographics analysis
- Property portfolio overview with charts
- Inheritance calculations for all living members
- Top inheritance opportunities
- Property valuations and recommendations

The generated PDF report includes professional charts, tables, and analysis suitable for sharing with family members or legal advisors.

## Push to Docker Hub

```bash
make push
```

## Persons Query

### Inheritance Confinements

```sql
SELECT d.id, d.name,
       att.url AS trello_url,
       a.storage_url
FROM inheritance_confinments d
LEFT JOIN LATERAL unnest(d.attachments) AS att(url) ON TRUE
LEFT JOIN attachments a
       ON att.url = a.trello_url
WHERE cardinality(d.attachments) > 0;
```
