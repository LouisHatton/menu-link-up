export function classnames(...names: string[]) {
	return names.join(' ');
}

export function validateEmail(email: string) {
	return String(email)
		.toLowerCase()
		.match(
			/^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
		);
}

export function validatePassword(password: string) {
	var re = /^(?=.*\d)(?=.*[!@#$%^&*])(?=.*[a-z])(?=.*[A-Z]).{8,}$/;
	return re.test(password);
}

export function sanitiseSlug(str: string) {
	return str
		.toLowerCase()
		.replace(/_/g, '') // Remove underscores
		.replace(/[^\w\s-]/g, '') // Remove non-word characters except spaces and hyphens
		.trim() // Trim leading/trailing spaces
		.replace(/\s+/g, '-') // Replace spaces with -
		.replace(/-+/g, '-'); // Replace multiple - with single -
}
