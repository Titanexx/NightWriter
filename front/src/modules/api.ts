import { ref } from 'vue';
import axios from 'axios';
import { Auth } from './auth';

export const useApiToken = (endpoint: string) =>{
	return useApi(endpoint,Auth.getToken());
}

export const useApi = (endpoint: string, token?: string) => {
	var headers = {}
	if(token){
		headers = {Authorization:  `Bearer ${token}`}
	} 

	const api = axios.create({
		baseURL: '/',
		headers: headers
	})
	const data = ref()
	const loading = ref(false)
	const error = ref()

	const post = (query?: Record<string, any>, config?: any) => {
		loading.value = true
		error.value = undefined

		return api.post(endpoint, query, config)
			.then(res => data.value = res.data)
			.catch(e => {
				error.value = e
			})
			.finally(() => loading.value = false)
            .then(() => {return data.value})
	}
    
	const put = (query?: Record<string, any>) => {
		loading.value = true
		error.value = undefined

		return api.put(endpoint, query)
			.then(res => data.value = res.data)
			.catch(e => {
				error.value = e
			})
			.finally(() => loading.value = false)
            .then(() => {return data.value})
	}

	const get = (query?: Record<string, any>) => {
		loading.value = true
		error.value = undefined

		return api.get(endpoint, query)
			.then(res => {data.value = res.data;})
			.catch(e => {
				error.value = e
			})
			.finally(() => loading.value = false)
            .then(() => {return data.value})
	}

    const syncGet = async (query?: Record<string, any>) => {
		loading.value = true
		error.value = undefined

		return await (api.get(endpoint, query)
			.then(res => {data.value = res.data;})
			.catch(e => {
				error.value = e
			})
			.finally(() => loading.value = false)
            .then(() => {return data.value}))
	}
    
	const del = (query?: Record<string, any>) => {
		loading.value = true
		error.value = undefined

		return api.delete(endpoint, query)
			.then(res => data.value = res.data)
			.catch(e => {
				error.value = e
			})
			.finally(() => loading.value = false)
            .then(() => {return data.value})
	}

	return {
		data, loading, error,
		get, post, put, del,
        syncGet
	}
}