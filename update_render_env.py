import requests
import json
import os

url = "https://api.render.com/v1/services/srv-d69gum9r0fns7383jbs0/env-vars"
api_key = "rnd_Dp2GTAkWZDy3OineypZt5TaWpCLQ"

headers = {
    "Authorization": f"Bearer {api_key}",
    "Content-Type": "application/json"
}

# Standard Pooler Format: postgres://postgres:[pw]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?pgbouncer=true
# Or even 5432 if they have IPv4 direct enabled, but 6543 is safer for pooler.
dsn = "postgresql://postgres:Inovar2026!Secure@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?pgbouncer=true"

payload = [
  {"key": "DATABASE_URL", "value": dsn},
  {"key": "SUPABASE_KEY", "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImJ4YnVwYm5qY2luZ2Z2anN6cmF1Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzAzODEyMjEsImV4cCI6MjA4NTk1NzIyMX0.OzBxS46bmR5OyxmS-DKFW7RRfEfVcgbhEKDWJSpMLOA"},
  {"key": "PORT", "value": "8080"},
  {"key": "JWT_SECRET", "value": "7f8e9d1a-2b3c-4d5e-6f7a-8b9c0d1e2f3a"},
  {"key": "SUPABASE_URL", "value": "https://bxbupbnjcingfvjszrau.supabase.co"}
]

response = requests.put(url, json=payload, headers=headers)

if response.status_code == 200:
    print("✅ Environment variables updated successfully on Render")
else:
    print(f"❌ Failed to update environment variables: {response.status_code}")
    print(response.text)
