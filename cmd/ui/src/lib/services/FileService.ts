import NetworkService from './NetworkService';

export type File = {
	id: string;
	name: string;
	userId: string;
	slug: string;
	createdAt: string;
	updatedAt: string;
	s3Bucket: string;
	s3Key: string;
};

export type NewFile = {
	name: string;
	slug: string;
};

export class FileService {
	listFiles(): Promise<File[]> {
		return NetworkService.get('/api/v1/files');
	}

	createFile(newFile: NewFile): Promise<File> {
		return NetworkService.post('/api/v1/files', newFile);
	}

	checkFileSlug(newFile: NewFile): Promise<boolean> {
		return NetworkService.post('/api/v1/files/check', newFile);
	}
}

export default new FileService();
