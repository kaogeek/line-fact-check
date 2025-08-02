import axios from 'axios';
import type { CountTopic, CountTopicCriteria, GetTopicCriteria, Topic } from '../type/topic';
import type { PaginationReq } from '../type/base';
import apiClient from '../client';
import { appendPaginationParams } from './base';

function buildTopicSearchParams(
  criteria: GetTopicCriteria | CountTopicCriteria,
  pagination?: PaginationReq
): URLSearchParams {
  const params = new URLSearchParams();

  if (criteria.codeLike) params.append('codeLike', criteria.codeLike);
  if (criteria.messageLike) params.append('like_message_text', criteria.messageLike);

  if ('statusIn' in criteria && criteria.statusIn?.length) {
    params.append('in_statuses', criteria.statusIn.join(','));
  }
  if ('idNotIn' in criteria && criteria.idNotIn?.length) {
    params.append('idNotIn', criteria.idNotIn.join(','));
  }

  if (pagination) {
    appendPaginationParams(params, pagination);
  }

  return params;
}

function patchData(topic: Topic): Topic {
  // TODO resolve this with real code
  return {
    ...topic,
    code: topic.id,
  };
}

export async function getTopics(criteria: GetTopicCriteria, pagination: PaginationReq): Promise<Topic[]> {
  const params = buildTopicSearchParams(criteria, pagination);
  const response = await apiClient.get<Topic[]>('/topics', { params });
  return response.data.map(patchData) ?? [];
}

export async function getTopicById(id: string): Promise<Topic | null> {
  try {
    const response = await apiClient.get<Topic>(`/topics/${id}`);
    return patchData(response.data);
  } catch (error) {
    if (axios.isAxiosError(error) && error.response?.status === 404) {
      return null;
    }
    throw error;
  }
}

export async function countTopics(criteria: CountTopicCriteria): Promise<CountTopic> {
  const params = buildTopicSearchParams(criteria);
  const response = await apiClient.get<CountTopic>('/topics/count', { params });
  return response.data;
}

export async function approveTopic(topicId: string): Promise<void> {
  await apiClient.patch(`/topics/${topicId}/approve`);
}

export async function rejectTopic(topicId: string): Promise<void> {
  await apiClient.patch(`/topics/${topicId}/reject`);
}
