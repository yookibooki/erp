#!/bin/bash

# Run migrations script for ERP SaaS
PGPASSWORD=erp_password psql -h localhost -U erp_user -d erp_saas -f ../migrations/01_initial_schema.sql