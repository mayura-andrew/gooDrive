#!/usr/bin/env python3
"""
gooDrive Analytics Tracker
Fetches and displays GitHub repository statistics
"""

import requests
import json
from datetime import datetime
from typing import Dict, Any

REPO_OWNER = "mayura-andrew"
REPO_NAME = "gooDrive"
API_BASE = f"https://api.github.com/repos/{REPO_OWNER}/{REPO_NAME}"

def fetch_repo_stats() -> Dict[str, Any]:
    """Fetch repository statistics from GitHub API"""
    try:
        response = requests.get(API_BASE)
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        print(f"Error fetching repo stats: {e}")
        return {}

def fetch_release_stats() -> Dict[str, Any]:
    """Fetch release download statistics"""
    try:
        response = requests.get(f"{API_BASE}/releases")
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        print(f"Error fetching release stats: {e}")
        return []

def calculate_downloads(releases: list) -> Dict[str, int]:
    """Calculate total and per-release downloads"""
    total_downloads = 0
    release_breakdown = {}
    
    for release in releases:
        release_downloads = 0
        for asset in release.get('assets', []):
            release_downloads += asset.get('download_count', 0)
        
        release_breakdown[release['tag_name']] = release_downloads
        total_downloads += release_downloads
    
    return {
        'total': total_downloads,
        'breakdown': release_breakdown
    }

def display_stats(repo_data: Dict[str, Any], downloads: Dict[str, Any]):
    """Display formatted statistics"""
    print("\n" + "=" * 60)
    print(f"📊 gooDrive Analytics Dashboard")
    print(f"📅 Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("=" * 60)
    
    # Repository Stats
    print("\n🌟 Repository Metrics:")
    print(f"   Stars:              {repo_data.get('stargazers_count', 0):,}")
    print(f"   Forks:              {repo_data.get('forks_count', 0):,}")
    print(f"   Watchers:           {repo_data.get('watchers_count', 0):,}")
    print(f"   Open Issues:        {repo_data.get('open_issues_count', 0):,}")
    
    # Download Stats
    print(f"\n📥 Download Statistics:")
    print(f"   Total Downloads:    {downloads['total']:,}")
    
    if downloads['breakdown']:
        print(f"\n   By Release:")
        for tag, count in downloads['breakdown'].items():
            print(f"      {tag:20} {count:,} downloads")
    
    # Project Info
    print(f"\n📦 Project Information:")
    print(f"   Language:           {repo_data.get('language', 'N/A')}")
    print(f"   Size:               {repo_data.get('size', 0):,} KB")
    print(f"   Created:            {repo_data.get('created_at', 'N/A')[:10]}")
    print(f"   Last Updated:       {repo_data.get('updated_at', 'N/A')[:10]}")
    print(f"   Default Branch:     {repo_data.get('default_branch', 'N/A')}")
    
    # URLs
    print(f"\n🔗 Quick Links:")
    print(f"   Repository:         {repo_data.get('html_url', 'N/A')}")
    print(f"   Stars:              https://github.com/{REPO_OWNER}/{REPO_NAME}/stargazers")
    print(f"   Issues:             https://github.com/{REPO_OWNER}/{REPO_NAME}/issues")
    print(f"   Releases:           https://github.com/{REPO_OWNER}/{REPO_NAME}/releases")
    
    print("\n" + "=" * 60)
    print("💡 Tip: Run this script weekly to track growth!")
    print("=" * 60 + "\n")

def save_to_csv(repo_data: Dict[str, Any], downloads: Dict[str, Any]):
    """Save statistics to CSV file for tracking over time"""
    import csv
    from pathlib import Path
    
    csv_file = Path(__file__).parent.parent / "analytics_history.csv"
    file_exists = csv_file.exists()
    
    with open(csv_file, 'a', newline='') as f:
        writer = csv.writer(f)
        
        # Write header if file is new
        if not file_exists:
            writer.writerow([
                'Date', 'Stars', 'Forks', 'Watchers', 
                'Open Issues', 'Total Downloads', 'Latest Release Downloads'
            ])
        
        # Get latest release downloads
        latest_downloads = 0
        if downloads['breakdown']:
            latest_release = list(downloads['breakdown'].keys())[0]
            latest_downloads = downloads['breakdown'][latest_release]
        
        # Write data
        writer.writerow([
            datetime.now().strftime('%Y-%m-%d'),
            repo_data.get('stargazers_count', 0),
            repo_data.get('forks_count', 0),
            repo_data.get('watchers_count', 0),
            repo_data.get('open_issues_count', 0),
            downloads['total'],
            latest_downloads
        ])
    
    print(f"✅ Data saved to {csv_file}")

def main():
    """Main function"""
    print("🔍 Fetching gooDrive analytics...")
    
    repo_data = fetch_repo_stats()
    releases = fetch_release_stats()
    downloads = calculate_downloads(releases)
    
    if repo_data:
        display_stats(repo_data, downloads)
        
        # Ask to save to CSV
        try:
            save_option = input("\n💾 Save to CSV for historical tracking? (y/n): ").lower()
            if save_option == 'y':
                save_to_csv(repo_data, downloads)
        except (KeyboardInterrupt, EOFError):
            print("\n\nSkipped saving to CSV.")
    else:
        print("❌ Failed to fetch repository data.")

if __name__ == "__main__":
    main()
