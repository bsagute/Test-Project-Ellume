import http from 'k6/http';
import { sleep } from 'k6';

//Used to Add Virtul Users & Duration in Script itself

// export let options = {
//   vus: 10,
//   i:10,
//   duration: '30s',
// };


// We can use below to do the load testing i.e ramp up & down your request in preticular time 
//Stress 
export let options = {
    stages: [
    { duration: '5s', target: 20 },
    { duration: '10s', target: 10 },
    { duration: '5s', target: 0 },
  ],
}

 export let count =0
 export let myCount=0
export default function () {
    myCount++
console.log("MyLogCount :- ",myCount);
    http.get('http://localhost:4000/GetCityListService');
    TestFunction()
    sleep(1);
//  http.get('http://localhost:4000/GetMasterURLs');
//     sleep(1);
//  http.get("'http://test.k6.io'");
//     sleep(1);

    
}

//    k6 run -i 10 --vus 10 --insecure-skip-tls-verify ./performance/basicTest.js


export function TestFunction() {
    count=count+1
    console.log("TEST IS CALLED ",count);
     http.get('http://localhost:4000/GetMasterURLs');
    sleep(1);
}

