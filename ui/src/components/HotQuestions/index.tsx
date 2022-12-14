import { FC, useEffect, useState } from 'react';
import { Card, ListGroup, ListGroupItem } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { Icon } from '@/components';
import { useHotQuestions } from '@/services';

const HotQuestions: FC = () => {
  const { t } = useTranslation('translation', { keyPrefix: 'question' });
  const [questions, setQuestions] = useState<any>([]);
  const { data: questionRes } = useHotQuestions();

  useEffect(() => {
    const questionResp = questionRes?.list;
    if (Array.isArray(questionResp)) {
      setQuestions(questionResp);
    }
  }, [questionRes]);

  return (
    <Card>
      <Card.Header className="text-nowrap text-capitalize">
        {t('hot_questions')}
      </Card.Header>
      <ListGroup variant="flush">
        {questions.map((li) => {
          return (
            <ListGroupItem
              key={li.id}
              as={Link}
              to={`/questions/${li.id}`}
              action>
              <div className="link-dark">{li.title}</div>
              {li.answer_count > 0 ? (
                <div
                  className={`d-flex align-items-center fs-14 mt-1 ${
                    li.accepted_answer_id > 0
                      ? 'link-success'
                      : 'link-secondary'
                  }`}>
                  {li.accepted_answer_id >= 1 ? (
                    <Icon name="check-circle-fill" />
                  ) : (
                    <Icon name="chat-square-text-fill" />
                  )}
                  <span className="ms-1">
                    {t('x_answers', { count: li.answer_count })}
                  </span>
                </div>
              ) : null}
            </ListGroupItem>
          );
        })}
      </ListGroup>
    </Card>
  );
};
export default HotQuestions;
