import { memo, FC } from 'react';
import { ListGroupItem } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import { Icon, Tag, FormatTime, BaseUserCard } from '@/components';
import type { SearchResItem } from '@/common/interface';
import { escapeRemove } from '@/utils';

interface Props {
  data: SearchResItem;
}
const Index: FC<Props> = ({ data }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'question' });
  if (!data?.object_type) {
    return null;
  }
  let itemUrl = `/questions/${data.object.id}`;
  if (data.object_type === 'answer') {
    itemUrl = `/questions/${data.object.question_id}/${data.object.id}`;
  }
  return (
    <ListGroupItem className="py-3 px-0">
      <div className="mb-2 clearfix">
        <span
          className="float-start me-2 badge text-bg-dark"
          style={{ marginTop: '2px' }}>
          {data.object_type === 'question' ? 'Q' : 'A'}
        </span>
        <a className="h5 mb-0 link-dark text-break" href={itemUrl}>
          {data.object.title}
          {data.object.status === 'closed'
            ? ` [${t('closed', { keyPrefix: 'question' })}]`
            : null}
        </a>
      </div>
      <div className="d-flex flex-wrap align-items-center fs-14 text-secondary mb-2">
        <BaseUserCard data={data.object?.user_info} showAvatar={false} />

        <span className="split-dot" />
        <FormatTime
          time={data.object?.created_at}
          className="me-3"
          preFix={data.object_type === 'question' ? 'asked' : 'answered'}
        />
        <div className="d-flex align-items-center my-2 my-sm-0">
          <div className="d-flex align-items-center me-3">
            <Icon name="hand-thumbs-up-fill me-1" />
            <span> {data.object?.vote_count}</span>
          </div>
          <div
            className={`d-flex align-items-center ${
              data.object?.accepted ? 'text-success' : ''
            }`}>
            {data.object?.accepted ? (
              <Icon name="check-circle-fill me-1" />
            ) : (
              <Icon name="chat-square-text-fill me-1" />
            )}
            <span>{data.object?.answer_count}</span>
          </div>
        </div>
      </div>

      {data.object?.excerpt && (
        <p className="fs-14 text-truncate-2 mb-2 last-p text-break">
          {escapeRemove(data.object.excerpt)}
        </p>
      )}

      {data.object?.tags?.map((item) => {
        return <Tag key={item.slug_name} className="me-1" data={item} />;
      })}
    </ListGroupItem>
  );
};

export default memo(Index);
