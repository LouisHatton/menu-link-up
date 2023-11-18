import { updateEmail, updateProfile, type User } from 'firebase/auth';
import AuthenticationService from './AuthenticationService';
import NetworkService from './NetworkService';

export type DbUser = {
	id: string;
	displayName: string;
	email: string;
	emailVerified: boolean;
	stripeCustomerId: string;
	providerId: string;
};

class UserService {
	async getUser(): Promise<DbUser> {
		let currentUser = AuthenticationService.currentUser;
		if (!currentUser) throw new Error('Not currently logged in');
		return NetworkService.get(`api/v1/users/${currentUser.uid}`);
	}

	async deleteUser(): Promise<unknown> {
		let currentUser = AuthenticationService.currentUser;
		if (!currentUser) throw new Error('Not currently logged in');
		return NetworkService.delete(`api/v1/users/${currentUser.uid}`);
	}

	async updateUserDisplayName(displayName: string) {
		let currentUser = AuthenticationService.currentUser;
		if (!currentUser) throw new Error('Not currently logged in');
		if (displayName !== currentUser?.displayName) {
			updateProfile(currentUser, { displayName });
		}
	}
}

export default new UserService();
