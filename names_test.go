package main

import (
	"encoding/json"
	"testing"
)

func TestGetAll(t *testing.T) {
	data :=
		`[
  "bardusco.id", 
  "bargain.id", 
  "bargok.id", 
  "barhaughogda.id", 
  "barheine.id", 
  "barista.id", 
  "barn.id", 
  "barnabee.id", 
  "barnaby.id", 
  "barnesdavid.id", 
  "barneysfarm.id", 
  "barometer.id", 
  "baronmichaelserovey.id", 
  "baronzero.id", 
  "barretlee.id", 
  "barritrad.id", 
  "barry.id", 
  "barry123.id", 
  "barry_j_brady.id", 
  "barryblaha.id", 
  "barrysilbert.id", 
  "barryteoh.id", 
  "barstool.id", 
  "barstoolsports.id", 
  "barthobartho.id", 
  "barthorstman.id", 
  "bartleeten.id", 
  "bartprokop.id", 
  "bartwilson.id", 
  "bartwr.id", 
  "barustore.id", 
  "bas.id", 
  "basecamp.id", 
  "basechain.id", 
  "basement.id", 
  "bash.id", 
  "bashco.id", 
  "basheer.id", 
  "basic.id", 
  "basichouse.id", 
  "basit121.id", 
  "basleenders.id", 
  "bass_pct.id", 
  "bassemj.id", 
  "bastian.id", 
  "bastienrobert.id", 
  "basudev.id", 
  "bata_pc888.id", 
  "batak15.id", 
  "batchy66.id", 
  "batghana.id", 
  "bathing.id", 
  "baths.id", 
  "batman11.id", 
  "batpeace.id", 
  "battalion999.id", 
  "battery.id", 
  "battles.id", 
  "bauerleather.id", 
  "baxterlf.id", 
  "bayard.id", 
  "bayer.id", 
  "bayern.id", 
  "baygail.id", 
  "bayite.id", 
  "bayless.id", 
  "baymax007.id", 
  "bazaamazon.id", 
  "bazaarassistant.id", 
  "bazaarbookstore.id", 
  "bazaarcasino.id", 
  "bazaarchic.id", 
  "bazaarcity.id", 
  "bazaarelectronics.id", 
  "bazaarfantasies.id", 
  "bazaarhost.id", 
  "bazaarinabox.id", 
  "bazaarjewelery.id", 
  "bazaarphones.id", 
  "bazaarprints.id", 
  "bazaarsales.id", 
  "bazaarslots.id", 
  "bazaarstores.id", 
  "bazmasta.id", 
  "bbabbrah.id", 
  "bballard1337.id", 
  "bbbr_design.id", 
  "bbc_new.id", 
  "bbdodd.id", 
  "bbh.id", 
  "bboru4911.id", 
  "bbowerman.id", 
  "bbrennan.id", 
  "bbrenner1.id", 
  "bbrown.id", 
  "bbva.id", 
  "bby203.id", 
  "bc4h10.id", 
  "bc_qondi77.id", 
  "bcarltonmicropower.id"
]`
	var schema []string
	err := json.Unmarshal([]byte(data), &schema)

	if err != nil {
		t.Fatal(err)
	}
}
