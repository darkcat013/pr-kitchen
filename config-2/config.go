package config

import "time"

// const DINING_HALL_URL = "http://localhost:8083/distribution"
const DINING_HALL_URL = "http://host.docker.internal:8083/distribution"
const PORT = "8082"
const LOGS_ENABLED = true

const TIME_UNIT = time.Millisecond * TIME_UNIT_COEFF
const TIME_UNIT_COEFF = 100

const MENU_PATH = "config/menu.json"
const COOKS_PATH = "config/cooks.json"
const APPARATUSES_PATH = "config/apparatuses.json"
