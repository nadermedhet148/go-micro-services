import http from 'k6/http'

export let options = {
  vus: 5,
  iterations: 100000,
};

export default function () {
  const regions = [
    "Egypt", "United States", "European Union", "Asia", "Africa",
    "Australia", "South America", "Canada", "India", "Russia",
    "China", "Brazil", "Mexico", "Japan", "South Korea",
    "South Africa", "Turkey", "France", "Germany", "Italy"
  ];
  const randomRegion = regions[Math.floor(Math.random() * regions.length)];

  var createWallet = http.post("http://localhost:8070/api/v1/wallets", JSON.stringify({
    name: "name",
    user_id: 1,
    region: randomRegion
  }));

  var rechargeWallet = http.post("http://localhost:8070/api/v1/wallets/recharge", JSON.stringify({
    wallet_id: createWallet.json().id,
    amount: 1000,
    region: randomRegion
  }))
}