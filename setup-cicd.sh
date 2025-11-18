#!/bin/bash
# Setup CI/CD with Cloud Build for go-chats project

set -e

PROJECT_ID="go-chats-478611"
REGION="us-central1"
REPO_NAME="go-chats"
GITHUB_OWNER="makxtr"

echo "🚀 Setting up CI/CD for go-chats..."

# 1. Enable Cloud Build API (если еще не включен)
echo "📦 Enabling Cloud Build API..."
gcloud services enable cloudbuild.googleapis.com --project=$PROJECT_ID

# 2. Grant Cloud Build permissions to deploy to Cloud Run
echo "🔐 Setting up permissions..."
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
  --role="roles/run.admin" \
  --condition=None

gcloud iam service-accounts add-iam-policy-binding \
  ${PROJECT_NUMBER}-compute@developer.gserviceaccount.com \
  --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser" \
  --project=$PROJECT_ID

echo ""
echo "✅ Permissions configured!"
echo ""
echo "📋 Next steps:"
echo ""
echo "1. Push cloudbuild.yaml files to your GitHub repo:"
echo "   git add auth/cloudbuild.yaml chat-server/cloudbuild.yaml"
echo "   git commit -m 'Add Cloud Build configuration'"
echo "   git push"
echo ""
echo "2. Connect GitHub repository to Cloud Build:"
echo "   https://console.cloud.google.com/cloud-build/triggers/connect?project=$PROJECT_ID"
echo ""
echo "3. Create triggers for each service:"
echo ""
echo "   Auth Service Trigger:"
echo "   - Name: deploy-auth-service"
echo "   - Event: Push to branch"
echo "   - Branch: ^main$"
echo "   - Configuration: auth/cloudbuild.yaml"
echo "   - Included files: auth/**"
echo ""
echo "   Chat Service Trigger:"
echo "   - Name: deploy-chat-service"
echo "   - Event: Push to branch"
echo "   - Branch: ^main$"
echo "   - Configuration: chat-server/cloudbuild.yaml"
echo "   - Included files: chat-server/**"
echo ""
echo "🎉 After setup, any push to main will automatically deploy!"
