The `tr-extractor` project extracts trello cards given a board ID, stores card metadata to a local database storage.

The project provides a user interface on [Notion](https://notion.com).

## Macro Architecture

- Backend
    - Golang deployed on Railway
- Database
    - Postgres (with Vector support) deployed in Railway 
- Automations
    - Make.com
- Frontend
    - Notion (rendedred on multiple views)
    - Google Sheet

## Tools

The following tools are used in this project:

| Tool            | Description                       | Fee |
|-----------------|-----------------------------------|----------|
| [Make.com](https://us2.make.com) | Automation Platform | $9 Monthly for 10,000 ops |
| [Railway](https://railway.com/) | App Deployment Platform: App + Postgres Database | $5 Monthly for 8 GB/8 vCPU |
| [Notion](https://notion.com) | Wiki, Databases, Sites, etc Platform | $10 Monthly |
| [Google Sheets](https://docs.google.com/spreadsheets) | Spreadsheet | Free Tier |

## Merge and Tag

- Assuming we have a working branch i.e. `my-branch`
  - `git add --all`
  - `git commit -am "Major stuff..."`
  - `git push`
  - `git checkout main`
  - `git merge my-branch`
  - `git tag -a v1.0.0 -m "my great work"`
  - `git tag` to make sure is is created.
  - `git push --tags` to push tags to Github.

## Deployment

Currently the deployment is manual to Railway. But the following are some improvements:

### Railway

- Install CLI
- Automate Deployment using API
- Export the database:

```bash
pg_dump -h monorail.proxy.rlwy.net -p 11397 -U postgres -d railway -F c -f ./dba/dumps/railway_backup_$(date +"%Y-%m-%d").dump
```

- Import the database:

```bash
pg_restore --host=monorail.proxy.rlwy.net --port=11397 --username=postgres --dbname=railway --format=c ./dba/dumps/backup_2025-02-13.dump
```

```bash
psql --host=monorail.proxy.rlwy.net --port=11397 --username=postgres --dbname=railway -f ./dba/dumps/backup_2025-02-13.sql
```

### Make.com

- Stop Automations via API
- Start Automations via API

## Automations

These automations require Trello board IDs and and a Trello API Key: 

| Automation      | Description                       | Interval | 
|-----------------|-----------------------------------|----------|
| Refresh            | Request property cards be pulled from Trello using API  | Every day at 6:00 AM |
| Google Sheet         | Upon the completion of the refresh, a webhook is triggered to run an automation to update Google sheet properties   | Triggered by Refresh |
| Notion         | Upon the completion of the refresh, a webhook is triggered to run an automation to update Notion properties database   | Triggered by Refresh |

## Phases

### Phase 1 - Organization and Population 

Phase 1 organizes Trello boards and populats them with useful and relevant data. These are the most importat boards:

- Properties
- Inheritance Confinments
- Supportive Documents
- Expenses

### Phase 2 - Automation 

- Pull Trello cards every 3 hours and store in a PostgresSQL database. 
- Query the Postgres database every 3 hours to update Google Sheet with multiple sheet names: one for each board type.
- Query the Postgres database every 3 hours to update Notion database and provide different views on the data.

### Phase 3 - Knowledge Base Powered by AI

- Run a RAG pipeline to produce a private knowledge base.
- Produce a front-end to query the RAG system.