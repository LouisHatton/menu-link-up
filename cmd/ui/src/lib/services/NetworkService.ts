import AuthenticationService from './AuthenticationService';

export type ApiError = {
	status: number;
	message: string;
};

async function handleNetworkError(r: Response): Promise<ApiError> {
	console.log(r);
	let json = await r.json();
	let message = 'unable to parse error';
	if ('status' in json) {
		message = json.status;
	}

	return {
		status: r.status,
		message
	};
}

class NetworkService {
	async fetch<T>(route: string, method = 'GET', body?: object) {
		let token = await AuthenticationService.getToken();
		if (!token) throw new Error('not logged in');

		let response = await fetch(route, {
			method,
			body: body ? JSON.stringify(body) : undefined,
			headers: {
				Authorization: 'Bearer ' + token,
				'content-type': 'application/json'
			}
		});

		if (!response.ok) {
			throw await handleNetworkError(response);
		}

		let json: unknown;
		try {
			json = await response.json();
		} catch {
			json = {};
		}

		return json as T;
	}

	async get<T>(route: string) {
		return this.fetch(route, 'GET') as T;
	}

	async post<T>(route: string, body: object) {
		return this.fetch(route, 'POST', body) as T;
	}

	async delete<T>(route: string) {
		return this.fetch(route, 'DELETE') as T;
	}
}

export default new NetworkService();
