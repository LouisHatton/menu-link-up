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

export type FileUpload = {
	url: string;
};

export class FileService {
	listFiles(): Promise<File[]> {
		return NetworkService.get('/api/v1/files');
	}

	createFile(newFile: NewFile): Promise<FileUpload> {
		return NetworkService.post('/api/v1/files', newFile);
	}

	uploadFile(url: string, file: BodyInit): Promise<Response> {
		return fetch(url, {
			method: 'PUT',
			body: file,
			headers: {
				'content-type': 'application/pdf'
			}
		});
	}

	deleteFile(fileId: string): Promise<{}> {
		return NetworkService.delete(`/api/v1/files/${fileId}`);
	}

	getFileLink(fileId: string): Promise<string> {
		return NetworkService.get(`/api/v1/files/${fileId}/link`);
	}

	checkFileSlug(newFile: NewFile): Promise<boolean> {
		return NetworkService.post('/api/v1/files/check', newFile);
	}
}

export default new FileService();
