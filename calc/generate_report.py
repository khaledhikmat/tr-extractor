#!/usr/bin/env python3
"""
Inheritance Report Generator

Generates a comprehensive PDF report analyzing the Alshamaa family inheritance system.
Includes demographics, property analysis, inheritance calculations, and recommendations.
"""

import json
import pandas as pd
import numpy as np
from datetime import datetime
from collections import defaultdict, Counter
import matplotlib.pyplot as plt
import seaborn as sns
from io import BytesIO
import base64

# PDF generation
from reportlab.pdfgen import canvas
from reportlab.lib.pagesizes import A4, letter
from reportlab.lib import colors
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import inch
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle, PageBreak
from reportlab.platypus import Image as RLImage
from reportlab.lib.enums import TA_CENTER, TA_LEFT, TA_RIGHT


class InheritanceReportGenerator:
    def __init__(self, persons_file="persons.json", properties_file="properties.json"):
        """Initialize the report generator with data files."""
        self.spouse_share_percentage = 0.08  # 8% spouse share from app.py
        self.total_shares = 2400  # Total shares system from app.py
        
        # Load data
        with open(persons_file, 'r', encoding='utf-8') as f:
            self.persons_data = json.load(f)
        
        with open(properties_file, 'r', encoding='utf-8') as f:
            self.properties_data = json.load(f)
        
        # Create DataFrames
        self.persons_df = pd.DataFrame(self.persons_data)
        self.properties_df = pd.DataFrame(self.properties_data)
        
        # Create person map for inheritance calculations
        self.person_map = {p['name']: p for p in self.persons_data}
        
        # Set up PDF styles
        self.styles = getSampleStyleSheet()
        self.title_style = ParagraphStyle('CustomTitle', parent=self.styles['Title'], 
                                         fontSize=24, spaceAfter=30, alignment=TA_CENTER)
        self.heading_style = ParagraphStyle('CustomHeading', parent=self.styles['Heading1'], 
                                          fontSize=16, spaceAfter=12, textColor=colors.darkblue)
        self.subheading_style = ParagraphStyle('CustomSubHeading', parent=self.styles['Heading2'], 
                                             fontSize=14, spaceAfter=10)
        
    # -------------------------
    # HELPER FUNCTIONS (from app.py)
    # -------------------------
    def get_children(self, person):
        return person.get('children', [])
    
    def get_spouses(self, person):
        return person.get('spouses', [])
    
    def is_deceased(self, person):
        return bool(person.get('death_year'))
    
    def calculate_inheritance_share(self, owner_name, target_name):
        """Calculate inheritance share using the algorithm from app.py"""
        def helper(current_name, current_share):
            if current_name == target_name:
                return current_share
            current_person = self.person_map.get(current_name)
            if not current_person:
                return 0
            
            children = self.get_children(current_person)
            spouses = self.get_spouses(current_person)
            deceased = self.is_deceased(current_person)

            if deceased:
                num_spouses = len(spouses)
                num_children = len(children)

                spouse_share = current_share * self.spouse_share_percentage if num_spouses > 0 else 0
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
    
    def calculate_inheritance(self, person, property_item):
        """Calculate inheritance value for a person from a property"""
        ownership_share = property_item.get('shares', 0) / self.total_shares
        family_inheritance_share = self.calculate_inheritance_share(property_item['owner'], person['name'])
        actual_inheritance_share = ownership_share * family_inheritance_share
        
        property_value = property_item['area'] * property_item['square_meter_price']
        inheritance_value = property_value * actual_inheritance_share
        
        return family_inheritance_share, inheritance_value
    
    # -------------------------
    # ANALYSIS FUNCTIONS
    # -------------------------
    def analyze_demographics(self):
        """Analyze family demographics."""
        analysis = {}
        
        # Basic counts
        total_persons = len(self.persons_data)
        living_count = len([p for p in self.persons_data if not self.is_deceased(p)])
        deceased_count = total_persons - living_count
        
        analysis['total_persons'] = total_persons
        analysis['living_count'] = living_count  
        analysis['deceased_count'] = deceased_count
        analysis['living_percentage'] = (living_count / total_persons) * 100
        
        # Gender distribution
        gender_counts = Counter([p.get('gender', 'Unknown') for p in self.persons_data])
        analysis['gender_distribution'] = dict(gender_counts)
        
        # Geographic distribution
        residence_countries = Counter([p.get('residence_country', 'Unknown') for p in self.persons_data if not self.is_deceased(p)])
        analysis['geographic_distribution'] = dict(residence_countries)
        
        # Education levels
        education_counts = Counter([p.get('education', 'Unknown') for p in self.persons_data])
        analysis['education_distribution'] = dict(education_counts)
        
        # Profession distribution
        profession_counts = Counter([p.get('profession', 'Unknown') for p in self.persons_data])
        analysis['profession_distribution'] = dict(profession_counts)
        
        return analysis
    
    def analyze_properties(self):
        """Analyze property portfolio."""
        analysis = {}
        
        # Basic property statistics
        total_properties = len(self.properties_data)
        total_area = sum([p['area'] for p in self.properties_data])
        total_value = sum([p['area'] * p['square_meter_price'] * (p['shares'] / self.total_shares) for p in self.properties_data])
        
        analysis['total_properties'] = total_properties
        analysis['total_area'] = total_area
        analysis['total_value'] = total_value
        analysis['average_property_size'] = total_area / total_properties
        analysis['average_property_value'] = total_value / total_properties
        
        # Property categories
        possessed_count = sum([1 for p in self.properties_data if p['possessed']])
        unsold_count = sum([1 for p in self.properties_data if p['unsold']])
        organized_count = sum([1 for p in self.properties_data if p['organized']])
        effects_count = sum([1 for p in self.properties_data if p['effects']])
        
        analysis['category_counts'] = {
            'possessed': possessed_count,
            'unsold': unsold_count,
            'organized': organized_count,
            'effects': effects_count
        }
        
        # Location distribution
        location_counts = Counter([p['location'] for p in self.properties_data])
        analysis['location_distribution'] = dict(location_counts)
        
        # Property owners
        owner_counts = Counter([p['owner'] for p in self.properties_data])
        analysis['owner_distribution'] = dict(owner_counts)
        
        # Top properties by value
        properties_with_values = []
        for p in self.properties_data:
            value = p['area'] * p['square_meter_price']
            properties_with_values.append({
                'name': p['name'],
                'owner': p['owner'],
                'area': p['area'],
                'value': value,
                'location': p['location']
            })
        
        properties_with_values.sort(key=lambda x: x['value'], reverse=True)
        analysis['top_properties'] = properties_with_values[:10]
        
        return analysis
    
    def analyze_inheritance_opportunities(self):
        """Analyze inheritance opportunities for living family members."""
        living_persons = [p for p in self.persons_data if not self.is_deceased(p)]
        inheritance_analysis = []
        
        for person in living_persons:
            # Parse name to get first and last name for sorting
            name_parts = person['name'].split()
            if len(name_parts) >= 2:
                first_name = ' '.join(name_parts[:-1])  # Everything except last word
                last_name = name_parts[-1]  # Last word
            else:
                first_name = person['name']
                last_name = ''
            
            person_inheritance = {
                'name': person['name'],
                'first_name': first_name,
                'last_name': last_name,
                'residence': f"{person.get('residence_city', '')}, {person.get('residence_country', '')}",
                'total_inheritance_value': 0,
                'property_count': 0,
                'best_properties': []
            }
            
            property_inheritances = []
            
            for prop in self.properties_data:
                if person['name'] != prop['owner']:  # Can't inherit from yourself
                    share, value = self.calculate_inheritance(person, prop)
                    if share > 0:
                        property_inheritances.append({
                            'property': prop['name'],
                            'owner': prop['owner'],
                            'location': prop['location'],
                            'share': share,
                            'value': value,
                            'property_total_value': prop['area'] * prop['square_meter_price']
                        })
            
            # Sort by inheritance value
            property_inheritances.sort(key=lambda x: x['value'], reverse=True)
            
            person_inheritance['total_inheritance_value'] = sum([p['value'] for p in property_inheritances])
            person_inheritance['property_count'] = len(property_inheritances)
            person_inheritance['best_properties'] = property_inheritances[:5]  # Top 5
            
            inheritance_analysis.append(person_inheritance)
        
        # Sort by last name, then first name (alphabetical)
        inheritance_analysis.sort(key=lambda x: (x['last_name'].lower(), x['first_name'].lower()))
        
        return inheritance_analysis
    
    # -------------------------
    # VISUALIZATION FUNCTIONS
    # -------------------------
    def create_pie_chart(self, data_dict, title, figsize=(8, 6)):
        """Create a pie chart and return as base64 encoded image."""
        fig, ax = plt.subplots(figsize=figsize)
        
        labels = list(data_dict.keys())
        sizes = list(data_dict.values())
        colors = plt.cm.Set3(np.linspace(0, 1, len(labels)))
        
        wedges, texts, autotexts = ax.pie(sizes, labels=labels, colors=colors, autopct='%1.1f%%', startangle=90)
        ax.set_title(title, fontsize=14, fontweight='bold', pad=20)
        
        # Make percentage text bold and larger
        for autotext in autotexts:
            autotext.set_color('white')
            autotext.set_fontsize(10)
            autotext.set_weight('bold')
        
        plt.tight_layout()
        
        # Save to BytesIO and encode as base64
        buffer = BytesIO()
        plt.savefig(buffer, format='png', dpi=300, bbox_inches='tight')
        buffer.seek(0)
        image_data = buffer.getvalue()
        buffer.close()
        plt.close()
        
        return base64.b64encode(image_data).decode()
    
    def create_bar_chart(self, data_dict, title, xlabel, ylabel, figsize=(10, 6)):
        """Create a bar chart and return as base64 encoded image."""
        fig, ax = plt.subplots(figsize=figsize)
        
        items = list(data_dict.items())
        items.sort(key=lambda x: x[1], reverse=True)  # Sort by value descending
        
        labels = [item[0] for item in items]
        values = [item[1] for item in items]
        
        bars = ax.bar(labels, values, color=plt.cm.Set2(np.linspace(0, 1, len(labels))))
        ax.set_title(title, fontsize=14, fontweight='bold', pad=20)
        ax.set_xlabel(xlabel)
        ax.set_ylabel(ylabel)
        
        # Add value labels on bars
        for bar in bars:
            height = bar.get_height()
            ax.annotate(f'{height:,.0f}' if height > 1000 else f'{height}',
                       xy=(bar.get_x() + bar.get_width() / 2, height),
                       xytext=(0, 3),  # 3 points vertical offset
                       textcoords="offset points",
                       ha='center', va='bottom',
                       fontsize=9, fontweight='bold')
        
        plt.xticks(rotation=45, ha='right')
        plt.tight_layout()
        
        # Save to BytesIO and encode as base64
        buffer = BytesIO()
        plt.savefig(buffer, format='png', dpi=300, bbox_inches='tight')
        buffer.seek(0)
        image_data = buffer.getvalue()
        buffer.close()
        plt.close()
        
        return base64.b64encode(image_data).decode()
    
    # -------------------------
    # PDF GENERATION
    # -------------------------
    def generate_report(self, output_filename="inheritance_report.pdf"):
        """Generate the complete PDF report."""
        print("Generating Inheritance Report...")
        
        # Perform analyses
        demographics = self.analyze_demographics()
        properties = self.analyze_properties()
        inheritance = self.analyze_inheritance_opportunities()
        
        # Create PDF document
        doc = SimpleDocTemplate(output_filename, pagesize=letter, 
                              rightMargin=72, leftMargin=72, 
                              topMargin=72, bottomMargin=18)
        
        # Container for PDF elements
        story = []
        
        # Title page
        story.append(Paragraph("Alshamaa Family Inheritance Report", self.title_style))
        story.append(Spacer(1, 20))
        story.append(Paragraph(f"Generated on: {datetime.now().strftime('%B %d, %Y')}", 
                              self.styles['Normal']))
        story.append(PageBreak())
        
        # Executive Summary
        story.append(Paragraph("Executive Summary", self.heading_style))
        story.append(Paragraph(f"""
        This comprehensive report analyzes the Alshamaa family inheritance system based on current family data and property portfolio. 
        The analysis covers {demographics['total_persons']} family members and {properties['total_properties']} properties with 
        a total estimated value of ${properties['total_value']:,.2f}.
        """, self.styles['Normal']))
        story.append(Spacer(1, 20))
        
        # Key Statistics Table
        key_stats_data = [
            ['Metric', 'Value'],
            ['Total Family Members', f"{demographics['total_persons']}"],
            ['Living Members', f"{demographics['living_count']} ({demographics['living_percentage']:.1f}%)"],
            ['Deceased Members', f"{demographics['deceased_count']}"],
            ['Total Properties', f"{properties['total_properties']}"],
            ['Total Land Area', f"{properties['total_area']:,.0f} sqm"],
            ['Total Portfolio Value', f"${properties['total_value']:,.2f}"],
            ['Average Property Value', f"${properties['average_property_value']:,.2f}"],
        ]
        
        stats_table = Table(key_stats_data, colWidths=[3*inch, 2*inch])
        stats_table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.grey),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, 0), 12),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('BACKGROUND', (0, 1), (-1, -1), colors.beige),
            ('GRID', (0, 0), (-1, -1), 1, colors.black)
        ]))
        story.append(stats_table)
        story.append(PageBreak())
        
        # Demographics Section
        story.append(Paragraph("Family Demographics Analysis", self.heading_style))
        
        # Living vs Deceased chart
        living_data = {'Living': demographics['living_count'], 'Deceased': demographics['deceased_count']}
        living_chart = self.create_pie_chart(living_data, "Family Status Distribution")
        living_image = RLImage(BytesIO(base64.b64decode(living_chart)), width=4*inch, height=3*inch)
        story.append(living_image)
        story.append(Spacer(1, 20))
        
        # Geographic distribution
        story.append(Paragraph("Geographic Distribution (Living Members)", self.subheading_style))
        geo_chart = self.create_bar_chart(demographics['geographic_distribution'], 
                                         "Living Family Members by Country", "Country", "Number of Members")
        geo_image = RLImage(BytesIO(base64.b64decode(geo_chart)), width=6*inch, height=3.5*inch)
        story.append(geo_image)
        story.append(PageBreak())
        
        # Property Analysis Section  
        story.append(Paragraph("Property Portfolio Analysis", self.heading_style))
        
        # Property categories chart
        categories_chart = self.create_pie_chart(properties['category_counts'], "Property Categories")
        categories_image = RLImage(BytesIO(base64.b64decode(categories_chart)), width=4*inch, height=3*inch)
        story.append(categories_image)
        story.append(Spacer(1, 20))
        
        # Location distribution
        story.append(Paragraph("Properties by Location", self.subheading_style))
        location_chart = self.create_bar_chart(properties['location_distribution'],
                                             "Properties by Location", "Location", "Number of Properties")
        location_image = RLImage(BytesIO(base64.b64decode(location_chart)), width=6*inch, height=3.5*inch)
        story.append(location_image)
        story.append(Spacer(1, 20))
        
        # All Properties Table
        story.append(Paragraph("List of All Properties", self.subheading_style))
        all_props_data = [['Property Name', 'Owner', 'Location', 'Price ($/sqm)', 'Area (sqm)', 'Shares', 'Effective Value ($)']]

        # Calculate effective value for each property: area * cost_per_sqm * shares/2400
        all_properties_with_values = []
        for prop in self.properties_data:
            effective_value = prop['area'] * prop['square_meter_price'] * (prop['shares'] / self.total_shares)
            all_properties_with_values.append({
                'name': prop['name'],
                'owner': prop['owner'],
                'location': prop['location'],
                'price': prop['square_meter_price'],
                'area': prop['area'],
                'shares': prop['shares'],
                'effective_value': effective_value
            })
        
        # Sort by area (descending - largest first)
        all_properties_with_values.sort(key=lambda x: x['area'], reverse=True)
        
        for prop in all_properties_with_values:
            all_props_data.append([
                prop['name'],
                prop['owner'][:20] + '...' if len(prop['owner']) > 20 else prop['owner'],
                prop['location'],
                f"{prop['price']:.2f}",
                f"{prop['area']:,.0f}",
                f"{prop['shares']:.2f}",
                f"${prop['effective_value']:,.0f}"
            ])
        
        all_props_table = Table(all_props_data, colWidths=[1.4*inch, 1.3*inch, 1*inch, 0.8*inch, 0.7*inch, 1.3*inch])
        all_props_table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.darkblue),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'CENTER'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, 0), 9),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('BACKGROUND', (0, 1), (-1, -1), colors.lightblue),
            ('GRID', (0, 0), (-1, -1), 1, colors.black),
            ('FONTSIZE', (0, 1), (-1, -1), 7),
        ]))
        story.append(all_props_table)
        story.append(PageBreak())
        
        # Inheritance Analysis Section
        story.append(Paragraph("Inheritance Analysis", self.heading_style))
        story.append(Paragraph("""
        This section analyzes potential inheritance opportunities for living family members based on the current 
        inheritance calculation algorithm. The algorithm considers family relationships, spouse shares (8%), 
        and proportional distribution among children.Inheritor names are sorted alphabetically by last name.
        """, self.styles['Normal']))
        story.append(Spacer(1, 20))
        
        # All Living Inheritors
        story.append(Paragraph("All Living Inheritors", self.subheading_style))
        inheritance_data = [['Name', 'Residence', 'Properties', 'Inheritance ($)']]
        for person in inheritance:  # Show all living inheritors, not just top 10
            inheritance_data.append([
                person['name'][:25] + '...' if len(person['name']) > 25 else person['name'],
                person['residence'][:20] + '...' if len(person['residence']) > 20 else person['residence'],
                str(person['property_count']),
                f"${person['total_inheritance_value']:,.0f}"
            ])
        
        inheritance_table = Table(inheritance_data, colWidths=[2*inch, 1.5*inch, 0.8*inch, 1.2*inch])
        inheritance_table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.darkgreen),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'CENTER'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('FONTSIZE', (0, 0), (-1, 0), 10),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('BACKGROUND', (0, 1), (-1, -1), colors.lightgreen),
            ('GRID', (0, 0), (-1, -1), 1, colors.black),
            ('FONTSIZE', (0, 1), (-1, -1), 8),
        ]))
        story.append(inheritance_table)
        story.append(Spacer(1, 20))
        
        story.append(PageBreak())
        
        # Recommendations Section
        story.append(Paragraph("Insights & Assumptions", self.heading_style))
        
        recommendations = [
            "1. **Property Management**: The Shabaa district properties represent the highest values and should be prioritized for proper documentation and management.",
            
            "2. **Geographic Diversification**: Family members are spread across multiple countries, which may complicate inheritance processes. Consider establishing clear legal frameworks in each jurisdiction.",
            
            f"3. **High-Value Assets**: The top property (Shabaa-151) alone is worth over ${properties['top_properties'][0]['value']:,.0f}. Special attention should be given to its inheritance planning.",
            
            "4. **Living Beneficiaries**: Focus on the inheritance planning for the 35 living family members, particularly those with significant inheritance potential.",
            
            "5. **Documentation**: Ensure all inheritance documents and legal frameworks are properly established, especially for properties marked as 'organized'.",
            
            "6. **Family Communication**: Regular family meetings should be held to discuss inheritance plans and ensure transparency among all stakeholders.",
            
            "7. **Price per Square Meter**: Estimated value at $100 per sqm for all properties. This may or may not be accurate.",
            
            "8. **Selected Properties**: The only properties that are considered for this report are those that have owner, area, shares and unsold.",
            
            "9. **Ownership**: It is assumed that the selected properties are fully owned by the designated owner and that no downstream claims exist."
        ]
        
        for rec in recommendations:
            story.append(Paragraph(rec, self.styles['Normal']))
            story.append(Spacer(1, 10))
        
        # Footer
        story.append(Spacer(1, 30))
        story.append(Paragraph(f"Report generated on {datetime.now().strftime('%B %d, %Y')} using the Alshamaa Family Inheritance Calculator", 
                              self.styles['Italic']))
        
        # Build PDF
        doc.build(story)
        print(f"Report generated successfully: {output_filename}")
        return output_filename


def main():
    """Main function to generate the inheritance report."""
    generator = InheritanceReportGenerator()
    output_file = generator.generate_report()
    print(f"\nInheritance report generated: {output_file}")
    print("The report includes:")
    print("- Family demographics analysis")
    print("- Property portfolio overview")
    print("- Inheritance calculations for all living members")
    print("- Recommendations for estate planning")


if __name__ == "__main__":
    main()