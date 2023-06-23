import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
    vus: 1,
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
    sleep(1);
}
