package config

import "time"

const DINING_HALL_URL = "http://localhost:8081/distribution"
const LOGS_ENABLED = true

const TIME_UNIT = time.Millisecond * TIME_UNIT_COEFF
const TIME_UNIT_COEFF = 100

const MENU_PATH = "config/menu.json"
const COOKS_PATH = "config/cooks.json"
const APPARATUSES_PATH = "config/apparatuses.json"
