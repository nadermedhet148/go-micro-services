import http from 'k6/http'

export let options = {
  vus: 10,
  iterations: 10000000,
};

export default function () {
  // var createWallet = http.post("http://localhost:8070/api/v1/wallets", JSON.stringify({
  //   name: "name",
  //   user_id : 1
  // }))

    var rechargeWallet = http.post("http://localhost:8070/api/v1/wallets/recharge", JSON.stringify({
    wallet_id: 1,
    amount: 1000,
  }))
  

}