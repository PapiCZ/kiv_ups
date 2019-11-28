#!/bin/bash

request_id() {
    echo -n "$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 5 | head -n 1)"
}

random_name() {
    echo -n "$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1)"
}

echo -n "|200|21|$(request_id)|{\"name\":\"$(random_name)\"}|201|40|$(request_id)|{\"name\":\"$(random_name)\",\"players_limit\":$((RANDOM % 30 + 10))}"
