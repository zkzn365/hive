import { ListGroup } from 'react-bootstrap';
import { Link } from 'react-router-dom';

import classNames from 'classnames';
import { isEmpty } from 'lodash';

import { FormatTime, Empty } from '@/components';

const Inbox = ({ data, handleReadNotification }) => {
  if (!data) {
    return null;
  }
  if (isEmpty(data)) {
    return <Empty />;
  }
  return (
    <ListGroup className="border-top border-bottom" variant="flush">
      {data.map((item) => {
        const { comment, question, answer } =
          item?.object_info?.object_map || {};
        let url = '';
        switch (item.object_info.object_type) {
          case 'question':
            url = `/questions/${item.object_info.object_id}`;
            break;
          case 'answer':
            url = `/questions/${question}/${item.object_info.object_id}`;
            break;
          case 'comment':
            url = `/questions/${question}/${answer}?commentId=${comment}`;
            break;
          default:
            url = '';
        }
        return (
          <ListGroup.Item
            key={item.id}
            className={classNames('py-3', !item.is_read && 'warning')}>
            <div>
              {item.user_info.status !== 'deleted' ? (
                <Link to={`/users/${item.user_info.username}`}>
                  {item.user_info.display_name}{' '}
                </Link>
              ) : (
                <span>{item.user_info.display_name} </span>
              )}
              {item.notification_action}{' '}
              <Link to={url} onClick={() => handleReadNotification(item.id)}>
                {item.object_info.title}
              </Link>
            </div>
            <div className="text-secondary">
              <FormatTime time={item.update_time} />
            </div>
          </ListGroup.Item>
        );
      })}
    </ListGroup>
  );
};

export default Inbox;
