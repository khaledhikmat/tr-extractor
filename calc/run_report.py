#!/usr/bin/env python3
"""
Simple script to run the inheritance report generator with different options.
"""

import sys
import os
from generate_report import InheritanceReportGenerator

def main():
    """Run the report generator with command line options."""
    
    # Parse command line arguments
    output_filename = "inheritance_report.pdf"
    if len(sys.argv) > 1:
        output_filename = sys.argv[1]
    
    # Check if data files exist
    if not os.path.exists("persons.json"):
        print("Error: persons.json file not found!")
        return
    
    if not os.path.exists("properties.json"):
        print("Error: properties.json file not found!")
        return
    
    try:
        # Create and run the report generator
        generator = InheritanceReportGenerator()
        output_file = generator.generate_report(output_filename)
        
        print(f"\n✅ SUCCESS: Inheritance report generated successfully!")
        print(f"📄 File: {output_file}")
        print(f"📊 Report includes:")
        print("   • Family demographics and statistics")
        print("   • Property portfolio analysis with charts")
        print("   • Inheritance calculations for all living members")
        print("   • Top inheritance opportunities")
        print("   • Property valuation and recommendations")
        print(f"\n💡 You can now share this PDF report with family members or advisors.")
        
    except Exception as e:
        print(f"❌ Error generating report: {str(e)}")
        print("Please check that all required packages are installed:")
        print("pip install pandas numpy matplotlib seaborn reportlab")

if __name__ == "__main__":
    main()