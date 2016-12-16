import * as mocha from "mocha";
import { assert } from "chai";

async function getGeoIP() {
    const res = await fetch("http://freegeoip.net/json/");
    const json = await res.json();
    return json as FreeGeoIP;
}

describe("fetch tests", () => {
    it("get ip", done => {
        getGeoIP().then(geoip => {
            // console.log(geoip);
            assert.isString(geoip.ip);
            assert.isString(geoip.country_code);
            done();
        });
    });
});