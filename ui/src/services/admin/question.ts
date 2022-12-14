import qs from 'qs';
import useSWR from 'swr';

import request from '@/utils/request';
import type * as Type from '@/common/interface';

export const useQuestionSearch = (params: Type.AdminContentsReq) => {
  const apiUrl = `/answer/admin/api/question/page?${qs.stringify(params)}`;
  const { data, error, mutate } = useSWR<Type.ListResult, Error>(
    [apiUrl],
    request.instance.get,
  );
  return {
    data,
    isLoading: !data && !error,
    error,
    mutate,
  };
};

export const changeQuestionStatus = (
  question_id: string,
  status: Type.AdminQuestionStatus,
) => {
  return request.put('/answer/admin/api/question/status', {
    question_id,
    status,
  });
};
