import os
from dotenv import load_dotenv
import streamlit as st
import pandas as pd
import json
import plotly.express as px
import boto3
from botocore.exceptions import ClientError
from datetime import datetime, timedelta

# -------------------------
# CONFIG
# -------------------------
USERS = {
    "admin": os.getenv("ADMIN_PASSWORD", "admin123"),  # Default password for admin
    "user": os.getenv("USER_PASSWORD", "user123")       # Default password for
}

SPOUSE_SHARE_PERCENTAGE = 0.08  # configurable spouse share
AWS_REGION = "us-east-2"        # your S3 region
PRESIGNED_EXPIRY = 3600         # seconds

# -------------------------
# AWS S3 CLIENT
# -------------------------
s3_client = boto3.client("s3", region_name=AWS_REGION)

def generate_presigned_url(s3_url):
    """
    Convert s3 URL (https://bucket.s3.region.amazonaws.com/key) to presigned URL
    """
    try:
        if s3_url and s3_url.startswith("https://"):
            # Parse URL more robustly to handle different regions
            url_without_protocol = s3_url.replace("https://", "")
            if ".s3." in url_without_protocol and ".amazonaws.com/" in url_without_protocol:
                bucket_and_region = url_without_protocol.split(".amazonaws.com/")[0]
                key = url_without_protocol.split(".amazonaws.com/")[1]
                bucket = bucket_and_region.split(".s3.")[0]
                
                # Extract region from URL
                region_part = bucket_and_region.split(".s3.")[1]
                if region_part:
                    # Create region-specific S3 client
                    region_s3_client = boto3.client("s3", region_name=region_part)
                else:
                    # Fallback to default client
                    region_s3_client = s3_client
                
                return region_s3_client.generate_presigned_url(
                    'get_object',
                    Params={'Bucket': bucket, 'Key': key},
                    ExpiresIn=PRESIGNED_EXPIRY
                )
        else:
            # If not a valid S3 URL, return the original URL
            return s3_url
    except (ClientError, IndexError, ValueError) as e:
        st.error(f"Error generating presigned URL for {s3_url}: {e}")
        # Return the original URL as fallback
        return s3_url
    return s3_url

# -------------------------
# LOAD DATA
# -------------------------
with open('persons.json') as f:
    persons_data = json.load(f)

with open('properties.json') as f:
    properties_data = json.load(f)

persons_df = pd.DataFrame(persons_data)
properties_df = pd.DataFrame(properties_data)

# -------------------------
# HELPER FUNCTIONS
# -------------------------
def build_person_map(persons):
    return {p['name']: p for p in persons}

person_map = build_person_map(persons_data)

def get_children(person):
    return person.get('children', [])

def get_spouses(person):
    return person.get('spouses', [])

def is_deceased(person):
    return bool(person.get('death_year'))

def calculate_inheritance_share(owner_name, target_name, person_map):
    def helper(current_name, current_share):
        if current_name == target_name:
            return current_share
        current_person = person_map.get(current_name)
        if not current_person:
            return 0
        
        children = get_children(current_person)
        spouses = get_spouses(current_person)
        deceased = is_deceased(current_person)

        if deceased:
            num_spouses = len(spouses)
            num_children = len(children)

            spouse_share = current_share * SPOUSE_SHARE_PERCENTAGE if num_spouses > 0 else 0
            children_share = current_share - spouse_share if num_children > 0 else 0

            if target_name in spouses and spouse_share > 0:
                return spouse_share / num_spouses

            if num_children > 0:
                share_per_child = children_share / num_children
                for child_name in children:
                    share = helper(child_name, share_per_child)
                    if share > 0:
                        return share
            return 0

        else:
            num_children = len(children)
            if num_children > 0:
                share_per_child = current_share / num_children
                for child_name in children:
                    share = helper(child_name, share_per_child)
                    if share > 0:
                        return share
            return 0

    return helper(owner_name, 1.0)

def calculate_inheritance(person, property, person_map):
    total_shares = 2400
    ownership_share = property.get('shares', 0) / total_shares
    family_inheritance_share = calculate_inheritance_share(property['owner'], person['name'], person_map)
    actual_inheritance_share = ownership_share * family_inheritance_share

    property_value = property['area'] * property['square_meter_price']
    inheritance_value = property_value * actual_inheritance_share

    return family_inheritance_share, inheritance_value

# -------------------------
# PAGES
# -------------------------
def dashboard():
    st.title("Dashboard")

    alive_count = len([p for p in persons_data if not p.get('death_year')])
    deceased_count = len(persons_data) - alive_count

    possessed_count = sum(1 for p in properties_data if p['possessed'])
    unsold_count = sum(1 for p in properties_data if p['unsold'])
    organized_count = sum(1 for p in properties_data if p['organized'])
    effects_count = sum(1 for p in properties_data if p['effects'])

    st.subheader("Persons Status")
    fig = px.pie(
        names=["Alive 🟢", "Deceased ⚰️"],
        values=[alive_count, deceased_count],
        color_discrete_sequence=px.colors.sequential.RdBu
    )
    st.plotly_chart(fig, use_container_width=True)

    st.subheader("Properties Categories")
    fig = px.pie(
        names=["Possessed 🏠", "Unsold 🏚️", "Organized 📦", "Effects 🎭"],
        values=[possessed_count, unsold_count, organized_count, effects_count],
        color_discrete_sequence=px.colors.sequential.Plasma
    )
    st.plotly_chart(fig, use_container_width=True)

def persons():
    st.title("👤 Persons")
    search_name = st.text_input("Search for a person")
    filtered_persons = [p for p in persons_data if search_name.lower() in p['name'].lower()]

    if not filtered_persons:
        st.write("No persons found.")
        return

    selected_name = st.selectbox("Select a person to see details", options=[p['name'] for p in filtered_persons])
    if selected_name:
        person = person_map[selected_name]
        st.markdown(f"## {person['name']}")

        st.markdown(f"- 👤 **Gender:** {person.get('gender', 'N/A')}")
        st.markdown(f"- 🎓 **Education:** {person.get('education', 'N/A')}")
        st.markdown(f"- 🧑🏻‍💻 **Profession:** {person.get('profession', 'N/A')}")
        st.markdown(f"- 🌍 **Birth:** {person.get('birth_city', 'N/A')}, {person.get('birth_country', 'N/A')} ")
        st.markdown(f"- 🛋️ **Residence:** {person.get('residence_city', 'N/A')}, {person.get('residence_country', 'N/A')} ")
        if person.get("death_year", "") != "":
            st.markdown(f"- ⚰️ **Death:** {person.get('death_day', 'N/A')}/{person.get('death_month','N/A')}/{person.get('death_year','N/A')} ({person.get('death_city','N/A')})")

        # include chilren and spouses only if deceased
        if person.get("death_year", "") != "":
            st.markdown("### 👶 Children")
            for child in person.get("children", []):
                st.markdown(f"- {child}")

            st.markdown("### 🧑‍🤝👰‍♀️ Spouses")
            for spouse in person.get("spouses", []):
                st.markdown(f"- {spouse }")

        st.markdown("### 📄 Documents")
        for doc_key, label in [("photo","📸 Photo"), ("inheritance_confinement","📄 Inheritance Confinement")]:
            if person.get(doc_key):
                url = generate_presigned_url(person[doc_key])
                if url:
                    st.markdown(f"- [{label}]({url})")

        st.markdown("### 🏡 Properties Owned")
        owned_props = [p for p in properties_data if p['owner'] == person['name']]
        if owned_props:
            df_owned = pd.DataFrame([{
                "Property": p['name'],
                "Location": f"{p['location']}, {p['city']}, {p['country']}",
                "Area": f"{p['area']} {p['area_unit']}",
                "Shares": f"{p['shares']}",
                "Value": f"${p['area']*p['square_meter_price']:,.2f}"
            } for p in owned_props])
            st.dataframe(df_owned, use_container_width=True)
        else:
            st.write("None")

        st.markdown("### 🤔 Potential Inheritance")
        potential_inheritance = []
        for prop in properties_data:
            if person['name'] != prop['owner']:
                share, value = calculate_inheritance(person, prop, person_map)
                if share > 0:
                    potential_inheritance.append({
                        "Property": prop['name'],
                        "Owner": prop['owner'],
                        "Share": f"{share:.2%}",
                        "Value": f"${value:,.2f}"
                    })
        if potential_inheritance:
            st.dataframe(pd.DataFrame(potential_inheritance), use_container_width=True)
        else:
            st.write("None")

def properties():
    st.title("🏡 Properties")
    search_property = st.text_input("Search for a property")
    filtered_properties = [p for p in properties_data if search_property.lower() in p['name'].lower()]

    if not filtered_properties:
        st.write("No properties found.")
        return

    selected_name = st.selectbox("Select a property to see details", options=[p['name'] for p in filtered_properties])
    if selected_name:
        prop = next(p for p in properties_data if p['name'] == selected_name)
        st.markdown(f"## {prop['name']} 🏠")

        st.markdown(f"- 📖 **Description:** {prop.get('description', 'N/A')}")
        st.markdown(f"- 📍 **Location:** {prop.get('location')}, {prop.get('city')}, {prop.get('country')}")
        st.markdown(f"- 📏 **Area:** {prop.get('area')} {prop.get('area_unit')}")
        st.markdown(f"- 🧳 **Shares:** {prop.get('shares')} of 2400 total shares")
        st.markdown(f"- 💰 **Price per m²:** ${prop.get('square_meter_price', 0)}")
        st.markdown(f"- 👤 **Owner:** {prop['owner']}")
        st.markdown(f"- 🏷️ **Attrs:** Possessed: {'✅' if prop['possessed'] else '❌'}, Unsold: {'✅' if prop['unsold'] else '❌'}, Organized: {'✅' if prop['organized'] else '❌'}, Effects: {'✅' if prop['effects'] else '❌'}")

        st.markdown("### 📄 Documents")
        for doc in prop.get("documents", []):
            st.markdown(f"- [{doc['type']}]({generate_presigned_url(doc['url'])})")

        st.markdown("### 👥 Eligible Persons for Inheritance")
        eligible = []
        for person in persons_data:
            share, value = calculate_inheritance(person, prop, person_map)
            if share > 0:
                eligible.append({
                    "Person": person['name'],
                    "Share": f"{share:.2%}",
                    "Value": f"${value:,.2f}"
                })
        if eligible:
            st.dataframe(pd.DataFrame(eligible), use_container_width=True)
        else:
            st.write("No one eligible.")

# -------------------------
# LOGIN PAGE
# -------------------------
def login():
    st.markdown("""
        <style>
        .login-box {
            max-width: 400px;
            margin: auto;
            padding: 2rem;
            border-radius: 10px;
            background-color: #f5f5f5;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        .login-box h2 {
            text-align: center;
            color: #333;
        }
        .login-box .stButton button {
            width: 100%;
            background-color: #4CAF50;
            color: white;
            font-size: 16px;
            height: 40px;
        }
        </style>
        <div class="login-box">
            <h2>Alshamaa Inheritance Calculator</h2>
        </div>
    """, unsafe_allow_html=True)

    username = st.text_input("Username", key="login_username")
    password = st.text_input("Password", type="password", key="login_password")
    login_button = st.button("Login")

    if login_button:
        if username in USERS and USERS[username] == password:
            st.session_state["logged_in"] = True
            st.session_state["username"] = username
            st.success(f"Welcome, {username}!")
            #st.experimental_rerun()
            st.stop()
        else:
            st.error("Invalid username or password")

# -------------------------
# MAIN
# -------------------------
def main():
    # Load .env
    dotenv_path = os.path.join(os.path.dirname(__file__), ".env")
    if os.path.exists(dotenv_path):
        load_dotenv(dotenv_path=dotenv_path)

    if "logged_in" not in st.session_state:
        st.session_state["logged_in"] = False

    if not st.session_state["logged_in"]:
        login()
    else:
        st.sidebar.title("🧭 Navigation")
        page = st.sidebar.radio("Choose a page", ["Dashboard", "Persons", "Properties"])
        if page == "Dashboard":
            dashboard()
        elif page == "Persons":
            persons()
        elif page == "Properties":
            properties()

        if st.sidebar.button("Logout"):
            st.session_state["logged_in"] = False
            #st.experimental_rerun()
            st.stop()

if __name__ == "__main__":
    main()
