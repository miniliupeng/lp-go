package main

import (
	"testing"
)

func TestUserMethodReceivers(t *testing.T) {
	u := User{ID: 1, Username: "Bob", Age: 25, IsActive: true}

	u.PrintGreeting()
	if u.Age != 25 {
		t.Fatalf("值接收者不应修改原对象")
	}

	u.CelebrateBirthday()
	if u.Age != 26 {
		t.Fatalf("指针接收者预期 Age=26，实际为 %d", u.Age)
	}
}

func TestEngineInterfacePolymorphism(t *testing.T) {
	gas := GasEngine{Gallons: 10, MPG: 30}     // 300
	elec := ElectricEngine{KWh: 50, MPKWh: 4} // 200

	if !CheckJourneyPlan(gas, 250) {
		t.Errorf("燃油车续航 300 应能到达 250")
	}
	if CheckJourneyPlan(elec, 250) {
		t.Errorf("电车续航 200 不应到达 250")
	}
}
