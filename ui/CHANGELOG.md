# Mainflux-UI Changelog

## 0.12.1 - 11. MAY 2021.

- NOISSUE - Fix nginx conf for groups
- NOISSUE - Add defaults.ini as docker volume for grafana container

## 0.12.0 - 07. APR 2021.

- NOISSUE - Fix menu style
- NOISSUE - Rm ng2-smart-table and npm audit fix
- NOISSUE - Update docker-compose for Mainflux 0.12.0
- UI-174 - Replace ng2-smart-table by the custom one in users and groups
- NOISSUE - Fix reader request
- NOISSUE - Fix .env for latest docker-compose version
- NOISSUE - Fix new disconnected readers filter
- UI-156 - Implement v, vb, vs and vd filters in message-monitor component
- UI-191 - Refactor twins to use generalized table component in states and definitions
- UI-185 - Add contrast to alternating table rows and generalize table template
- NOISSUE - Remap diff value fields to a field value in messages response
- UI-187 - Add string, bool and data types in message-monitor and http client
- NOISSUE - Replace ngIf with ngSwitch in table template
- UI-147 - Refactor channels/things connection dropdown menu
- UI-174 - Use custom table in opc-ua page
- NOISSUE - Remove experimental build and use only dev and prod build
- NOISSUE - Add separate pagination component to Gateway and Lora svcs
- NOISSUE - Split Table component into Table and Pagination components
- NOISSUE - Fix time unit in charts (use milliseconds)
- UI-174 - Use custom table in twins states page
- NOISSUE - Fix docker-compose for twins service
- UI-174 - Implement custom table and message-monitor in gateways page
- UI-174 - Implement custom table and message-monitor in lora page
- NOISSUE - Fix style of sub-menu items
- UI-171 - Change Menu and Header colors
- NOISSUE - Sort connections list by name ascendent
- NOISSUE - Add prefix and suffix inputs to configure reader path
- UI-141 - Replace Things, Channels and Twins ng2-smart-table by the custom one
- UI-155 - Implement table pagination with limit
- UI-165 - Create custom table component with pagination
- NOISSUE - separate nginx entrypoint for ui and main nginx
- UI-152 - Add messages reader monitor (json, table and chart modes)
- NOISSUE - Update InfluxDB and Postgres docker-compose images
- UI-157 - Update docker-compose to use auth service
- UI-152 - Implement datepicker filter in messages-table
- NOISSUE - Add interval service
- UI-143 - Add message table component and message value pipe
- UI-140 - Add type edit/create to table and remove type undefined
- NOISSUE - Add missing attributes to editor attr list
- NOISSUE - Add Dockerfile.experimental
- NOISSUE - Order Things and Channels by name ascendent
- NOISSUE - Add UTF-8 BOM to CSV file
- NOISSUE - Add customizable strings file
- NOISSUE - Add experimental conditional build
- NOISSUE - Add environment substitution for ui
- NOISSUE - Add details to setup instructions
- NOISSUE - Add messages sufix envar
- NOISSUE - Fix typo
- NOISSUE - Rename Dashboard -> Home
- NOISSUE - Fix error messages and xterm bug in gateways page
- NOISSUE - Rename Orgs -> Groups
- NOISSUE - Add Organisations (Users groups) page
- NOISSUE - Use new menu hierarchy and rename Devices->Things
- NOISSUE - Add Users page
- UI-121 - Upgrade Dockerfile Node.js to 14.14.0
- NOISSUE - Fix typo and remove commented code
- NOISSUE - Bump ckeditor from 4.7.3 to 4.12.1
- NOISSUE - Update akveo/ngx-admin from starter-kit branch
- NOISSUE - Configure charts to display other values
- NOISSUE - Improve twin definition editor ergonomy
- NOISSUE - Remove Loraserver, Tracing and Grafana pages
- NOISSUE - Use new colors from styles/mainflux.css
- NOISSUE - Set password minLength to 8
- NOISSUE - Make angular environment variables configurable through docker environment vars
- NOISSUE - Add user admin envars in docker-compose
- NOISSUE - Add new Dashboard
- NOISSUE - Fix ng2-charts version
- NOISSUE - Add provision service to UI composition
- NOISSUE - Fix typo in authorization.js
- NOISSUE - Use not connected Things and Channels endpoint
- NOISSUE - Add clean command to Makefile
- NOISSUE - Fixed typo in ssl Makefile

## 0.11.0 - 17. JUN 2020.

- NOISSUE - update vars and docker-compose
- NOISSUE - Renaming mac, gw password
- NOISSUE - Use log level error for VermeMQ docker
- UI-73 - Check if ID is defined before opening the table details
- NOISSUE - fix var naming use snake case
- Update some deprecated dependencies, fix missing dependencies and minimizing severity vulnerabilities
- NOISSUE - Sort charts data by message name
- UI-83 - Remove db reader logout on forbidden error
- Update nginx-x509.conf
- UI-79 - Check appID mapping before to create LoRa Devices
- UI-78 - Use snake_case for Lora and OPC-UA metadata fields
- NOISSUE - fix term variabl, and mqtt url
- NOISSUE - Update docker-compose and .env
- NOISSUE - Add type in Devices and Channels tables
- NOISSUE - Update docker-compose
- NOISSUE: adding remote terminal
- UI-32 - Verify if MAC exist for Gateways creation and edit
- UI-49 - Fix responsive buttons
- NOISSUE - Add delta field to twin's definition
- UI-23 - Implement Save table for all pages
- NOISSUE - Fix opcua name editing
- NOISSUE - Fix typos and mv rxjs imports to common.module
- NOISSUE - Fix linter errors
- NOISSUE - Check if subscription exist before creation
- NOISSUE - Add namespace and identifier to opc-ua browser
- Update README.md
- Update README.md
- UI-48 - Fix login sequence and add LoginComponent
- NOISSUE - Add SenML value types to current state message composer
- NOISSUE - Polish Twins lists
- NOISSUE - Add Twins to docker-compose and nginx conf
- NOISSUE - Enable edit twin name in twins table
- NOISSUE - Add store and search bar for browsed opc-ua server data
- NOISSUE - Export config
- NOISSUE - Refactor definition attribute display in history and editor
- NOISSUE - Add publishing mqtt to service subtopic
- NOISSUE - Replace ng2-smart-table search bars by Mainflux implementation
- NOISSUE - Add definition history page
- NOISSUE - Fix non-implemented save buttons
- [UI-33] MQTT url configurable from environment
- NOISSUE - URL encode subtopic query param
- NOISSUE - Fix get Things with metadata request
- NOISSUE - Add bulk upload in opc-ua page
- NOISSUE - Twins without thing
- UI-23 - Add bulk upload for Things and Channels
- NOISSUE - Add export configuration
- NOISSUE- Bump tinymce from 4.5.7 to 4.9.7
- NOISSUE - Improve messages service to fetch device or channel messages
- NOISSUE - Twins page
- NOISSUE - Add details to browsed OPC-UA nodes
- NOISSUE - Add password endpoint to nginx conf
- NOISSUE - add release command to Makefile
- NOISSUE - Add complete docker-compose + nginx configuration
- NOISSUE - Fix error catch and statusText
- NOISSUE - Check if the opcua node channel exist before creation
- NOISSUE - Simplify version urls
- NOISSUE - Add Browse card inOPC-UA page

## 0.10.0 - 02. DEC 2019.
### Features
- Dashboard: User informations.
- Things: system management.
- LoRa: Route map configuration for LoRa Server.
- Edge: MFX-1 Gateway control over the MFX-edged.
- Admin: DB, Tracing,  Admin


### Summary
https://github.com/mainflux/mainflux/milestone/9?closed=1
