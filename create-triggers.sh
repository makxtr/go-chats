#!/bin/bash
# Create Cloud Build triggers for go-chats project
# Run this AFTER connecting GitHub repo to Cloud Build

set -e

PROJECT_ID="go-chats-478611"
REGION="us-central1"
REPO_NAME="go-chats"
GITHUB_OWNER="makxtr"

echo "🚀 Creating Cloud Build triggers..."
echo ""
echo "⚠️  Make sure you've connected GitHub repo first at:"
echo "   https://console.cloud.google.com/cloud-build/triggers/connect?project=$PROJECT_ID"
echo ""
read -p "Have you connected the GitHub repo? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Please connect the repo first, then run this script again."
    exit 1
fi

# Get the repository resource name
echo "📋 Looking for connected repository..."

# List all repositories to find the one we need
REPO_FULL_NAME="${GITHUB_OWNER}_${REPO_NAME}"

echo ""
echo "Creating trigger for auth-service..."

# Create trigger for auth service
gcloud builds triggers create github \
  --name="deploy-auth-service" \
  --description="Deploy auth-service on push to main" \
  --repo-name="$REPO_NAME" \
  --repo-owner="$GITHUB_OWNER" \
  --branch-pattern="^main$" \
  --build-config="auth/cloudbuild.yaml" \
  --included-files="auth/**" \
  --region="$REGION" \
  --project="$PROJECT_ID"

echo "✅ Auth service trigger created!"
echo ""
echo "Creating trigger for chat-service..."

# Create trigger for chat service
gcloud builds triggers create github \
  --name="deploy-chat-service" \
  --description="Deploy chat-service on push to main" \
  --repo-name="$REPO_NAME" \
  --repo-owner="$GITHUB_OWNER" \
  --branch-pattern="^main$" \
  --build-config="chat-server/cloudbuild.yaml" \
  --included-files="chat-server/**" \
  --region="$REGION" \
  --project="$PROJECT_ID"

echo "✅ Chat service trigger created!"
echo ""
echo "🎉 All triggers created successfully!"
echo ""
echo "View triggers at:"
echo "https://console.cloud.google.com/cloud-build/triggers?project=$PROJECT_ID"
echo ""
echo "Now you can test by pushing changes to main branch:"
echo "  git add ."
echo "  git commit -m 'Test auto-deploy'"
echo "  git push"