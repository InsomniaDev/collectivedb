import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
    vus: 3,
    duration: '30s',
};

function randomString() {
    let val = (Math.random() + 1).toString(36).substring(7);
    return {
        key: val,
        data: "test"
    }
}

export default function () {
    let rando = randomString()
    http.post('http://192.168.1.19:31048/update', JSON.stringify(rando))
    // for (let i = 0; i < 10; i++) {
    http.get('http://192.168.1.19:31048/get/' + rando.key)

    rando.data = "new data"
    http.post('http://192.168.1.19:31048/update', JSON.stringify(rando))

    http.get('http://192.168.1.19:31048/get/' + rando.key)

    // }
    http.post('http://192.168.1.19:31048/delete', JSON.stringify(rando))
    // sleep(1);
}
