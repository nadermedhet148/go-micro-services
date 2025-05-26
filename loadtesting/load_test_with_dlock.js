import http from 'k6/http'

export let options = {
  vus: 10,
  iterations: 100,
};

export default function () {
  const regions = ["EG", "US", "EU", "AS", "AF"];
  const randomRegion = regions[Math.floor(Math.random() * regions.length)];

  // var createWallet = http.post("http://localhost:8070/api/v1/wallets", JSON.stringify({
  //   name: "name",
  //   user_id: 1,
  //   region: randomRegion
  // }));

    var rechargeWallet = http.post("http://localhost:8070/api/v1/wallets/recharge", JSON.stringify({
    wallet_id: 1,
    amount: 1000,
  }))
}