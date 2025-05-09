import http from 'k6/http'

export let options = {
  vus: 1000,
  iterations: 1000,
};

export default function () {
  const slotId = Math.floor(Math.random() * 100) + 1;
  let res = http.post("http://localhost:8001/api/v1/Ticket/d-lock", JSON.stringify({
    "slot_id": slotId,
    "ref_number": "test"
  }))

}