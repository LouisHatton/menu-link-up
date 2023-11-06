import NetworkService from './NetworkService';

export type Project = {
	id: string;
	name: string;
	slug: string;
	metadata: ProjectMetadata;
	users: string[];
	adminUsers: string[];
	roUsers: string[];
	rwUsers: string[];
};

export type ProjectMetadata = {
	createdBy: string;
	createdAt: string;
};

export type NewProject = {
	name: string;
	slug: string;
};

export class ProjectService {
	listProjects(): Promise<Project[]> {
		return NetworkService.get('/api/v1/projects');
	}

	createProject(newProject: NewProject): Promise<Project[]> {
		return NetworkService.post('/api/v1/projects', newProject);
	}

	checkProjectCreation(newProject: NewProject): Promise<boolean> {
		return NetworkService.post('/api/v1/projects/check', newProject);
	}
}

export default new ProjectService();
