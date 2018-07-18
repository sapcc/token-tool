#!/bin/bash

DOMAIN=${DOMAIN:-monsoon3}
PROJECT=${PROJECT:-consulting_dev}
PROJECT_DOMAIN=${PROJECT_DOMAIN:-$DOMAIN}
USERNAME=${USERNAME:-$USER}
KEYSTONE_ENDPOINT=${KEYSTONE_ENDPOINT:-https://identity-3.staging.cloud.sap/v3}
FORMAT=text

function HELP {
  cat << EOF
usage: $0
  -e KEYSTONE_ENDPOINT   (Default: $KEYSTONE_ENDPOINT)
  -u USERNAME            (Default: $PROJECT)
  -d USER_DOMAIN_NAME    (Default: $DOMAIN)
  -p PROJECT             (Default: $PROJECT)
  -q PROJECT_DOMAIN_NAME (Default: \$USER_DOMAIN_NAME)
  -f FORMAT              text,json,curlrc (Default: text)
EOF
}

while getopts u:d:p:e:q:f:vh FLAG; do
  case $FLAG in
    u)
      USERNAME=$OPTARG
      ;;
    d)
      DOMAIN=$OPTARG
      ;;
    p)
      PROJECT=$OPTARG
      ;;
    e)
      KEYSTONE_ENDPOINT=$OPTARG
      ;;
    q)
      PROJECT_DOMAIN=$OPTARG
      ;;
    f)
      FORMAT=$OPTARG
      ;;
    h)
      HELP
      exit 0
      ;;
    \?)
      HELP
      exit 1
      ;;
  esac
done

PASSWORD=$(security find-generic-password -a $USERNAME -s openstack -w 2>/dev/null)

if [ -z "$PASSWORD" ]; then
  read -p "Enter password for user $USERNAME: " -s PASSWORD
  echo
fi

payload=$(cat <<EOP
{ "auth": {
    "identity": {
      "methods": ["password"],
      "password": {
        "user": {
          "name": "$USERNAME",
          "domain": { "name": "$DOMAIN"   },
          "password": "$PASSWORD"
        }
      }
    },
    "scope": {
      "project": {
        "name": "$PROJECT",
        "domain": {"name":"$PROJECT_DOMAIN"}
      }
    }
  }
}
EOP
)

output=$(curl --silent -D /dev/stderr \
      -H "Content-Type: application/json" \
      -d "$payload" \
  $KEYSTONE_ENDPOINT/auth/tokens?nocatalog \
  2>&1)

# echo "==========="
# echo $output
# echo "==========="

json=$(echo "$output" |tail -1)
token=$( echo "$output" | sed -n 's/[Xx]-[Ss]ubject-[Tt]oken: \(.*\)/\1/p' | tr -d '\r')

if echo $json |grep -q error; then
  echo $json | jq -r .error.message
  exit 1
fi

case $FORMAT in
  json)
    echo $json | jq .
    ;;
  curlrc)
    echo "header \"X-Auth-Token: $token\""
    echo "header \"Content-Type: application/json\""
    ;;
  *)
  echo "Token: $token"
  echo -n "User: "
  echo $json | jq -r ".token.user.id"
  echo -n "Project: "
  echo $json | jq -r ".token.project.id"
  echo -n "Project Domain: "
  echo $json | jq -r ".token.project.domain.id"
  echo -n "Roles: "
  echo $json | jq -r '.token.roles |map(.name)|join(", ")'
esac