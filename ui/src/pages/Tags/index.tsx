import { useState } from 'react';
import { Container, Row, Col, Card, Button, Form } from 'react-bootstrap';
import { useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { Tag, Pagination, PageTitle, QueryGroup } from '@/components';
import { formatCount } from '@/utils';
import { useQueryTags, following } from '@/services';

const sortBtns = ['popular', 'name', 'newest'];

const Tags = () => {
  const [urlSearch] = useSearchParams();
  const { t } = useTranslation('translation', { keyPrefix: 'tags' });
  const [searchTag, setSearchTag] = useState('');

  const page = Number(urlSearch.get('page')) || 1;
  const sort = urlSearch.get('sort');

  const pageSize = 20;
  const { data: tags, mutate } = useQueryTags({
    page,
    page_size: pageSize,
    ...(searchTag ? { slug_name: searchTag } : {}),
    ...(sort ? { query_cond: sort } : {}),
  });

  const handleChange = (e) => {
    setSearchTag(e.target.value);
  };

  const handleFollow = (tag) => {
    following({
      object_id: tag.tag_id,
      is_cancel: tag.is_follower,
    }).then(() => {
      mutate();
    });
  };
  return (
    <>
      <PageTitle title={t('tags', { keyPrefix: 'page_title' })} />
      <Container className="py-3 my-3">
        <Row className="mb-4 d-flex justify-content-center">
          <Col xxl={10} sm={12}>
            <h3 className="mb-4">{t('title')}</h3>
            <div className="d-flex justify-content-between align-items-center flex-wrap">
              <Form>
                <Form.Group controlId="formBasicEmail">
                  <Form.Control
                    value={searchTag}
                    placeholder={t('search_placeholder')}
                    type="text"
                    onChange={handleChange}
                    size="sm"
                  />
                </Form.Group>
              </Form>
              <QueryGroup
                data={sortBtns}
                currentSort={sort || 'popular'}
                sortKey="sort"
                i18nKeyPrefix="tags.sort_buttons"
              />
            </div>
          </Col>

          <Col className="mt-4" xxl={10} sm={12}>
            <Row>
              {tags?.list?.map((tag) => (
                <Col
                  key={tag.slug_name}
                  xs={12}
                  lg={3}
                  md={4}
                  sm={6}
                  className="mb-4">
                  <Card className="h-100">
                    <Card.Body className="d-flex flex-column align-items-start">
                      <Tag className="mb-3" data={tag} />

                      <p className="fs-14 flex-fill text-break text-wrap text-truncate-4">
                        {tag.original_text}
                      </p>
                      <div className="d-flex align-items-center">
                        <Button
                          className={`me-2 ${tag.is_follower ? 'active' : ''}`}
                          variant="outline-primary"
                          size="sm"
                          onClick={() => handleFollow(tag)}>
                          {tag.is_follower
                            ? t('button_following')
                            : t('button_follow')}
                        </Button>
                        <span className="text-secondary fs-14 text-nowrap">
                          {formatCount(tag.question_count)} {t('tag_label')}
                        </span>
                      </div>
                    </Card.Body>
                  </Card>
                </Col>
              ))}
            </Row>
            <div className="d-flex justify-content-center">
              <Pagination
                currentPage={page}
                totalSize={tags?.count || 0}
                pageSize={pageSize}
              />
            </div>
          </Col>
        </Row>
      </Container>
    </>
  );
};

export default Tags;
