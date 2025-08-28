# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Streamlit-based family inheritance tracking application that manages genealogical data and property inheritance calculations. The app uses a simple authentication system and visualizes inheritance relationships through interactive charts and tables.

## Development Setup

### Environment Setup
```bash
python3 -m venv venv
source venv/bin/activate 
pip3 install -r requirements.txt
```

### Running the Application
```bash
streamlit run app.py
```

## Architecture

### Core Components

**Main Application (`app.py`)**
- Single-file Streamlit application with page-based navigation
- Authentication system with hardcoded credentials in `USERS` dict
- AWS S3 integration for document storage with presigned URL generation
- Three main pages: Dashboard, Persons, Properties

**Data Structure**
- `persons.json`: Family member data including birth/death info, relationships, and document URLs
- `properties.json`: Property data with ownership, location, and valuation information
- Relationships defined through `children` and `spouses` arrays in person objects

**Inheritance Calculation Engine**
- `calculate_inheritance_share()`: Recursive function calculating inheritance percentages
- `SPOUSE_SHARE_PERCENTAGE`: Configurable spouse inheritance ratio (default 8%)
- Algorithm handles deceased vs living persons differently
- Share calculation considers property ownership shares and family relationships

### Key Functions

**Person Management**
- `build_person_map()`: Creates name-to-person lookup dictionary
- `get_children()`, `get_spouses()`, `is_deceased()`: Relationship and status helpers

**AWS Integration**
- `generate_presigned_url()`: Converts S3 URLs to time-limited access URLs
- Configured for `us-east-2` region with 3600 second expiry

**Data Visualization**
- Uses Plotly Express for pie charts showing person/property status distributions
- Pandas DataFrames for tabular data display

## Configuration

### Constants in `app.py`
- `USERS`: Authentication credentials (line 12-15)
- `SPOUSE_SHARE_PERCENTAGE`: Inheritance calculation parameter (line 17)
- `AWS_REGION`: S3 region setting (line 18)
- `PRESIGNED_EXPIRY`: URL expiration time in seconds (line 19)

### Data Schema
- Person objects require `name`, optional `children`/`spouses` arrays, and death status
- Property objects require `name`, `owner`, `area`, `square_meter_price`, and categorical flags
- Property shares are calculated against total of 2400 shares

## Important Notes

- No test framework is currently implemented
- Application stores credentials in plain text (development only)
- S3 integration requires proper AWS credentials in environment
- JSON data files are loaded synchronously on startup
- Session state manages authentication and page navigation
- **Ignore the `ignore/` folder**: Contains old application versions that should not be modified