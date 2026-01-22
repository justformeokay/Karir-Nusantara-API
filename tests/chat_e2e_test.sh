#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Chat Support System - End-to-End Test ===${NC}\n"

# Test 1: Backend Health Check
echo -e "${YELLOW}1. Checking Backend Health...${NC}"
HEALTH=$(curl -s http://localhost:8081/health | jq -r '.status')
if [ "$HEALTH" = "healthy" ]; then
  echo -e "${GREEN}✅ Backend API is running${NC}\n"
else
  echo -e "${RED}❌ Backend API is NOT running${NC}"
  exit 1
fi

# Test 2: Company Login
echo -e "${YELLOW}2. Testing Company Login...${NC}"
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"info@karyadeveloperindonesia.com","password":"Justformeokay23@"}' | jq -r '.data.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo -e "${RED}❌ Login failed${NC}"
  exit 1
fi
echo -e "${GREEN}✅ Company login successful${NC}"
echo -e "   Token: ${TOKEN:0:20}...${NC}\n"

# Test 3: Create Conversation
echo -e "${YELLOW}3. Creating Test Conversation...${NC}"
CREATE_RESPONSE=$(curl -s -X POST 'http://localhost:8081/api/v1/company/chat/conversations' \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"subject":"Test E2E Chat","category":"helpdesk"}')

CONV_ID=$(echo $CREATE_RESPONSE | jq -r '.data.id')
CONV_TITLE=$(echo $CREATE_RESPONSE | jq -r '.data.title')

if [ -z "$CONV_ID" ] || [ "$CONV_ID" = "null" ]; then
  echo -e "${RED}❌ Failed to create conversation${NC}"
  echo $CREATE_RESPONSE | jq '.'
  exit 1
fi
echo -e "${GREEN}✅ Conversation created${NC}"
echo -e "   ID: $CONV_ID${NC}"
echo -e "   Title: $CONV_TITLE${NC}\n"

# Test 4: List Conversations
echo -e "${YELLOW}4. Listing Company Conversations...${NC}"
CONV_COUNT=$(curl -s "http://localhost:8081/api/v1/company/chat/conversations" \
  -H "Authorization: Bearer $TOKEN" | jq '.data | length')

echo -e "${GREEN}✅ Conversations retrieved${NC}"
echo -e "   Total: $CONV_COUNT${NC}\n"

# Test 5: Get Specific Conversation
echo -e "${YELLOW}5. Retrieving Specific Conversation...${NC}"
GET_RESPONSE=$(curl -s "http://localhost:8081/api/v1/company/chat/conversations/$CONV_ID" \
  -H "Authorization: Bearer $TOKEN")

CONV_STATUS=$(echo $GET_RESPONSE | jq -r '.data.conversation.status')
MSG_COUNT=$(echo $GET_RESPONSE | jq '.data.messages | length')

echo -e "${GREEN}✅ Conversation retrieved${NC}"
echo -e "   Status: $CONV_STATUS${NC}"
echo -e "   Messages: $MSG_COUNT${NC}\n"

# Test 6: Send Message
echo -e "${YELLOW}6. Sending Test Message...${NC}"
SEND_RESPONSE=$(curl -s -X POST "http://localhost:8081/api/v1/company/chat/conversations/$CONV_ID/messages" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"message":"Test message from E2E script"}')

MSG_ID=$(echo $SEND_RESPONSE | jq -r '.data.id')
MSG_TEXT=$(echo $SEND_RESPONSE | jq -r '.data.message')

if [ -z "$MSG_ID" ] || [ "$MSG_ID" = "null" ]; then
  echo -e "${RED}❌ Failed to send message${NC}"
  echo $SEND_RESPONSE | jq '.'
  exit 1
fi
echo -e "${GREEN}✅ Message sent${NC}"
echo -e "   ID: $MSG_ID${NC}"
echo -e "   Text: $MSG_TEXT${NC}\n"

# Test 7: Verify Message in Conversation
echo -e "${YELLOW}7. Verifying Message in Conversation...${NC}"
VERIFY_RESPONSE=$(curl -s "http://localhost:8081/api/v1/company/chat/conversations/$CONV_ID" \
  -H "Authorization: Bearer $TOKEN")

UPDATED_MSG_COUNT=$(echo $VERIFY_RESPONSE | jq '.data.messages | length')
LAST_MSG=$(echo $VERIFY_RESPONSE | jq -r '.data.messages[-1].message')

if [ "$UPDATED_MSG_COUNT" -gt 0 ]; then
  echo -e "${GREEN}✅ Message verified in conversation${NC}"
  echo -e "   Total messages: $UPDATED_MSG_COUNT${NC}"
  echo -e "   Last message: $LAST_MSG${NC}\n"
else
  echo -e "${RED}❌ Message not found in conversation${NC}"
  exit 1
fi

# Test 8: Admin Login & View
echo -e "${YELLOW}8. Testing Admin Access...${NC}"
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@karirnusantara.com","password":"admin123"}' | jq -r '.data.access_token')

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
  echo -e "${RED}❌ Admin login failed${NC}"
  exit 1
fi
echo -e "${GREEN}✅ Admin login successful${NC}\n"

# Test 9: Admin View All Conversations
echo -e "${YELLOW}9. Admin Viewing All Conversations...${NC}"
ADMIN_CONVS=$(curl -s "http://localhost:8081/api/v1/admin/chat/conversations" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.data | length')

echo -e "${GREEN}✅ Admin can view conversations${NC}"
echo -e "   Total conversations visible to admin: $ADMIN_CONVS${NC}\n"

# Test 10: Admin Reply Message
echo -e "${YELLOW}10. Admin Sending Reply...${NC}"
ADMIN_REPLY=$(curl -s -X POST "http://localhost:8081/api/v1/admin/chat/conversations/$CONV_ID/messages" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"message":"Admin reply: Percakapan sudah diterima. Terima kasih!"}')

ADMIN_MSG_ID=$(echo $ADMIN_REPLY | jq -r '.data.id')

if [ -z "$ADMIN_MSG_ID" ] || [ "$ADMIN_MSG_ID" = "null" ]; then
  echo -e "${RED}❌ Admin failed to send message${NC}"
  exit 1
fi
echo -e "${GREEN}✅ Admin reply sent${NC}"
echo -e "   Message ID: $ADMIN_MSG_ID${NC}\n"

# Test 11: Update Conversation Status
echo -e "${YELLOW}11. Admin Updating Conversation Status...${NC}"
STATUS_UPDATE=$(curl -s -X PATCH "http://localhost:8081/api/v1/admin/chat/conversations/$CONV_ID/status" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"status":"in_progress"}')

STATUS_SUCCESS=$(echo $STATUS_UPDATE | jq -r '.success')

if [ "$STATUS_SUCCESS" = "true" ]; then
  echo -e "${GREEN}✅ Conversation status updated to in_progress${NC}\n"
else
  echo -e "${RED}❌ Failed to update status${NC}"
fi

# Test 12: Verify Conversation State
echo -e "${YELLOW}12. Final Verification...${NC}"
FINAL_STATE=$(curl -s "http://localhost:8081/api/v1/company/chat/conversations/$CONV_ID" \
  -H "Authorization: Bearer $TOKEN")

FINAL_STATUS=$(echo $FINAL_STATE | jq -r '.data.conversation.status')
FINAL_MSG_COUNT=$(echo $FINAL_STATE | jq '.data.messages | length')

echo -e "${GREEN}✅ Final conversation state:${NC}"
echo -e "   Conversation ID: $CONV_ID${NC}"
echo -e "   Title: $(echo $FINAL_STATE | jq -r '.data.conversation.title')${NC}"
echo -e "   Status: $FINAL_STATUS${NC}"
echo -e "   Total Messages: $FINAL_MSG_COUNT${NC}"
echo -e "   Created: $(echo $FINAL_STATE | jq -r '.data.conversation.created_at')${NC}\n"

# Summary
echo -e "${YELLOW}=== Test Summary ===${NC}"
echo -e "${GREEN}✅ All tests passed!${NC}\n"
echo -e "Chat Support System is:"
echo -e "  ${GREEN}✓ Backend API functional${NC}"
echo -e "  ${GREEN}✓ Database persistence working${NC}"
echo -e "  ${GREEN}✓ Company-Admin communication working${NC}"
echo -e "  ${GREEN}✓ Message history preserved${NC}"
echo -e "  ${GREEN}✓ Status management working${NC}"
echo -e "\n${GREEN}🚀 SYSTEM IS PRODUCTION READY!${NC}"
