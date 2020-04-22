#!/bin/sh
docker tag golang-ubuntu1804-build:latest golang-ubuntu1804-build:backup
docker commit golang-ubuntu1804-build-jenkins golang-ubuntu1804-build:latest
