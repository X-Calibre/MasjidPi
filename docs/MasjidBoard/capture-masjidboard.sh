#!/bin/bash

BASE='https://api.masjidboardlive.com/mblapi?id='

curl -L "${BASE}1lTEEzl7sefO4W72c9iKxcUkB1ZReHhtZt9DFmCFfT_0" \
  -o fawkner-masjid.json

curl -L "${BASE}1ZK8NtqROdU3Ww4THcHkHyDJN2gu98HC1ovBbGO7iooY" \
  -o zakariyya-park-duzak.json

curl -L "${BASE}1asEQ0Ju83TPqBFHw7NbBAihAxMt5JQ2bJkbaWnwKf7k" \
  -o brits-jamia.json

curl -L "${BASE}170sRYVcxfOC-l3IGK0FeqWX8D8rM35afyl7nyL2SWHI" \
  -o brits-taqwa.json

curl -L "${BASE}1Zpg5LKfd_ZoEQsA0rsyWNBrUgY6QVaHnGdPfuKHF24A" \
  -o azaadville-darul-uloom.json
