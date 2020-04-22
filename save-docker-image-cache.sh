#!/bin/sh
docker tag golang-ubuntu-build:latest golang-ubuntu-build:backup
docker commit golang-ubuntu-build-jenkins golang-ubuntu-build:latest
