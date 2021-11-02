package main

type APIDataProvider interface {
    GetSpendingInLastHour()
    GetSpendingInRange(startHour, endHour int64)
}