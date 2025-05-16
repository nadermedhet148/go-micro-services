import http from 'k6/http'

export let options = {
  // vus: 1000,
  // iterations: 1000,
};

export default function () {
  const slotId = Math.floor(Math.random() * 100) + 1;
  // var createWallet = http.post("http://localhost:8080/api/v1/wallets", JSON.stringify({
  //   name: "name",
  //   user_id : 1
  // }))

    var rechargeWallet = http.post("http://localhost:8080/api/v1/wallets/recharge", JSON.stringify({
    wallet_id: 1,
    amount: 1000,
  }))

  console.log("createWallet", createWallet);
  

}