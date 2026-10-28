#!/bin/bash
printenv
# replace quad4Cloud domain and LDAP password in keycloak realm import
mkdir /opt/keycloak/data/import
sed -e "s/cloud.quad4.test/${OC_DOMAIN}/g" -e "s/ldap-admin-password/${LDAP_ADMIN_PASSWORD:-admin}/g" /opt/keycloak/data/import-dist/quad4-realm.json > /opt/keycloak/data/import/quad4-realm.json

# run original docker-entrypoint
/opt/keycloak/bin/kc.sh "$@"
