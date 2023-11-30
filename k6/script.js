import http from 'k6/http';
import {check} from 'k6';
import {Rate} from 'k6/metrics';

const reqRate = new Rate('http_req_rate');

export const options = {
    scenarios: {
        load: {
            executor: 'ramping-arrival-rate',
            startRate: 1,
            timeUnit: '1s',
            preAllocatedVUs: 20,
            stages: [
                {target: 20, duration: '40s'},
                {target: 0, duration: '20s'},
            ],
        },
    },
    thresholds: {
        'checks': ['rate>0.9'],
        'http_req_duration': ['p(95)<1000'],
        'http_req_rate{deployment:stable}': ['rate>=0'],
        'http_req_rate{deployment:canary}': ['rate>=0'],
    },
};

export default function () {
    const params = {
        headers: {
            'Host': 'sample.app',
            'Content-Type': 'application/json',
        },
    };

    const res = http.get(`http://localhost/success`, params);
    check(res, {
        'status code is 200': (r) => r.status === 200,
        'node is kind-control-plane': (r) => r.json().node === 'kind-control-plane',
        'namespace is sample-app': (r) => r.json().namespace === 'sample-app',
        'pod is sample-app-*': (r) => r.json().pod.includes('sample-app-'),
        'deployment is stable or canary': (r) => r.json().deployment === 'stable' || r.json().deployment === 'canary',
    });

    switch (res.json().deployment) {
        case 'stable':
            reqRate.add(true, { deployment: 'stable' });
            reqRate.add(false, { deployment: 'canary' });
            break;
        case 'canary':
            reqRate.add(false, { deployment: 'stable' });
            reqRate.add(true, { deployment: 'canary' });
            break;
    }
}
