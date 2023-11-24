import { updateEmail, updateProfile, type User } from 'firebase/auth';
import AuthenticationService from './AuthenticationService';
import NetworkService from './NetworkService';

export type DbUser = {
	id: string;
	email: string;
	stripeCustomerId: string;
	subscriptionStatus: 'trialing' | 'active' | 'cancelled';
	trialEnd: string | null;
};

export type Billing = {
	planName: string;
	billingInterval: 'month' | 'year';
	price: number;
	currentPeriodEnd?: string;
	cancelAtPeriodEnd: boolean;
	defaultPayment?: BillingDefaultPayment;
};

export type BillingDefaultPayment = {
	brand?: string;
	expiresMonth?: number;
	expiresYear?: number;
	lastFour?: string;
};

class UserService {
	async getUser(): Promise<DbUser> {
		let currentUser = AuthenticationService.currentUser;
		if (!currentUser) throw new Error('Not currently logged in');
		return NetworkService.get(`api/v1/users/${currentUser.uid}`);
	}

	async getUserBilling(): Promise<Billing> {
		let currentUser = AuthenticationService.currentUser;
		if (!currentUser) throw new Error('Not currently logged in');
		return NetworkService.get(`api/v1/users/${currentUser.uid}/billing`);
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
