#!/bin/bash

# Setup GCP Secret Manager for Neon Database URLs
set -e

PROJECT_ID="${1}"

if [ -z "$PROJECT_ID" ]; then
  echo "Usage: ./setup-gcp-secrets.sh YOUR_PROJECT_ID"
  exit 1
fi

echo "Setting up secrets for project: $PROJECT_ID"

# Check if neon-credentials.txt exists
if [ ! -f "neon-credentials.txt" ]; then
  echo "Error: neon-credentials.txt not found!"
  echo "Please create it with your Neon connection strings:"
  echo ""
  echo "AUTH_DATABASE_URL=postgresql://..."
  echo "CHAT_DATABASE_URL=postgresql://..."
  exit 1
fi

# Load credentials
export $(cat neon-credentials.txt | xargs)

if [ -z "$AUTH_DATABASE_URL" ] || [ -z "$CHAT_DATABASE_URL" ]; then
  echo "Error: AUTH_DATABASE_URL or CHAT_DATABASE_URL not set in neon-credentials.txt"
  echo "AUTH_DATABASE_URL: ${AUTH_DATABASE_URL:-not set}"
  echo "CHAT_DATABASE_URL: ${CHAT_DATABASE_URL:-not set}"
  exit 1
fi

echo "================================================"
echo "Enabling Secret Manager API..."
echo "================================================"
gcloud services enable secretmanager.googleapis.com --project=$PROJECT_ID

echo "================================================"
echo "Creating secrets..."
echo "================================================"

# Create or update auth database URL secret
if gcloud secrets describe auth-database-url --project=$PROJECT_ID &>/dev/null; then
  echo "Updating auth-database-url secret..."
  echo -n "$AUTH_DATABASE_URL" | gcloud secrets versions add auth-database-url \
    --data-file=- \
    --project=$PROJECT_ID
else
  echo "Creating auth-database-url secret..."
  echo -n "$AUTH_DATABASE_URL" | gcloud secrets create auth-database-url \
    --data-file=- \
    --replication-policy="automatic" \
    --project=$PROJECT_ID
fi

# Create or update chat database URL secret
if gcloud secrets describe chat-database-url --project=$PROJECT_ID &>/dev/null; then
  echo "Updating chat-database-url secret..."
  echo -n "$CHAT_DATABASE_URL" | gcloud secrets versions add chat-database-url \
    --data-file=- \
    --project=$PROJECT_ID
else
  echo "Creating chat-database-url secret..."
  echo -n "$CHAT_DATABASE_URL" | gcloud secrets create chat-database-url \
    --data-file=- \
    --replication-policy="automatic" \
    --project=$PROJECT_ID
fi

echo "================================================"
echo "Setting up permissions..."
echo "================================================"

# Get Cloud Build service account
# Cloud Build uses the Project Number, not Project ID, for its default service account
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"
COMPUTE_SA="${PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

# Grant permissions to Cloud Build
gcloud secrets add-iam-policy-binding auth-database-url \
  --member="serviceAccount:$BUILD_SA" \
  --role="roles/secretmanager.secretAccessor" \
  --project=$PROJECT_ID

gcloud secrets add-iam-policy-binding chat-database-url \
  --member="serviceAccount:$BUILD_SA" \
  --role="roles/secretmanager.secretAccessor" \
  --project=$PROJECT_ID

# Grant permissions to Cloud Run (Compute Engine default SA)
if [ ! -z "$COMPUTE_SA" ]; then
  gcloud secrets add-iam-policy-binding auth-database-url \
    --member="serviceAccount:$COMPUTE_SA" \
    --role="roles/secretmanager.secretAccessor" \
    --project=$PROJECT_ID

  gcloud secrets add-iam-policy-binding chat-database-url \
    --member="serviceAccount:$COMPUTE_SA" \
    --role="roles/secretmanager.secretAccessor" \
    --project=$PROJECT_ID
fi

echo "================================================"
echo "Setup complete! ✓"
echo "================================================"
echo ""
echo "Secrets created:"
echo "  - auth-database-url"
echo "  - chat-database-url"
echo ""
echo "Next steps:"
echo "1. Update cloudbuild.yaml files (already done if you run update script)"
echo "2. Push changes to trigger Cloud Build"
echo "3. Migrations will run automatically before deployment"
echo "================================================"
