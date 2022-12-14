import { FC, memo } from 'react';
import { ListGroup, ListGroupItem } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import { FormatTime } from '@/components';

interface Props {
  visible: boolean;
  data;
}

const Index: FC<Props> = ({ visible, data }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'personal' });
  if (!visible || !data?.length) {
    return null;
  }
  return (
    <ListGroup variant="flush">
      {data.map((item) => {
        return (
          <ListGroupItem className="d-flex py-3 px-0" key={item.object_id}>
            <div
              className={`me-3 text-end ${
                item.reputation > 0 ? 'text-success' : 'text-danger'
              }`}
              style={{ width: '40px' }}>
              {item.reputation > 0 ? '+' : ''}
              {item.reputation}
            </div>
            <div>
              <a
                className="text-break"
                href={`/questions/${
                  item.object_type === 'question'
                    ? item.object_id
                    : `${item.question_id}/${item.object_id}`
                }`}>
                {item.title}
              </a>
              <div className="d-flex align-items-center fs-14 text-secondary">
                <span>{item.reputation > 0 ? t('upvote') : t('downvote')}</span>
                <span className="split-dot" />
                <FormatTime time={item.created_at} className="me-4" />
              </div>
            </div>
          </ListGroupItem>
        );
      })}
    </ListGroup>
  );
};

export default memo(Index);
