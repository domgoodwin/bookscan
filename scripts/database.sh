#!/bin/bash

export POSTGRES_USER="postgres"
export POSTGRES_DB="postgres"
export POSTGRES_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@192.168.0.240:5432/${POSTGRES_DB}?sslmode=disable"