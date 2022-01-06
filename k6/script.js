import http from 'k6/http';
import {check, sleep} from 'k6';
import { Rate } from 'k6/metrics';

const reqRate = new Rate('http_req_rate');

export const options = {
    stages: [
        {target: 20, duration: '20s'},
        {target: 20, duration: '20s'},
        {target: 0, duration: '20s'},
    ],
    thresholds: {
        'checks': ['rate>0.9'],
        'http_req_duration': ['p(95)<1000'],
        'http_req_rate{deployment:stable}': ['rate>=0'],
        'http_req_rate{deployment:canary}': ['rate>=0'],
    },
};

export default function () {
    const url = `http://${__ENV.URL}`;
    const params = {
        headers: {
            'Host': 'sample.app',
            'Content-Type': 'application/json',
        },
    };

    const res = http.get(`${url}/success`, params)
    check(res, {
        'status code is 200': (r) => r.status === 200,
        'node is minikube': (r) => r.json().node === 'minikube',
        'namespace is sample-app': (r) => r.json().namespace === 'sample-app',
        'pod is sample-app-*': (r) => r.json().pod.includes('sample-app-'),
        'deployment is stable or canary': (r) => r.json().deployment === 'stable' || r.json().deployment === 'canary',
    });

    switch (res.json().deployment) {
        case 'stable':
            reqRate.add(true, { deployment: 'stable' })
            reqRate.add(false, { deployment: 'canary' })
            break
        case 'canary':
            reqRate.add(false, { deployment: 'stable' })
            reqRate.add(true, { deployment: 'canary' })
            break
    }

    sleep(1)
}
