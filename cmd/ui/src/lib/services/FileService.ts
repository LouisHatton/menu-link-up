import NetworkService, { type ApiError } from './NetworkService';

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
	fileSize: number;
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

	uploadFile(url: string, file: BodyInit, size: number): Promise<Response> {
		return fetch(url, {
			method: 'PUT',
			body: file,
			headers: {
				'content-type': 'application/pdf',
				'content-length': size.toString()
			}
		});
	}

	deleteFile(fileId: string): Promise<{}> {
		return NetworkService.delete(`/api/v1/files/${fileId}`);
	}

	getFileLink(fileId: string): Promise<string> {
		return NetworkService.get(`/api/v1/files/${fileId}/link`);
	}

	async checkFileSlug(slug: string): Promise<boolean> {
		try {
			await NetworkService.post('/api/v1/check-file', {
				slug
			});
			return true;
		} catch (err: unknown) {
			console.log(err as ApiError);
			return false;
		}
	}
}

export default new FileService();
