import { FC } from 'react';
import { Form, Table, Dropdown } from 'react-bootstrap';
import { useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import classNames from 'classnames';

import {
  Pagination,
  FormatTime,
  BaseUserCard,
  Empty,
  QueryGroup,
  Icon,
} from '@/components';
import * as Type from '@/common/interface';
import { useChangeModal, useChangeUserRoleModal, useToast } from '@/hooks';
import { useQueryUsers } from '@/services';
import { loggedUserInfoStore } from '@/stores';

import '../index.scss';

const UserFilterKeys: Type.UserFilterBy[] = [
  'all',
  // 'staff',
  'inactive',
  'suspended',
  'deleted',
];

const bgMap = {
  normal: 'text-bg-success',
  suspended: 'text-bg-danger',
  deleted: 'text-bg-danger',
  inactive: 'text-bg-secondary',
};

const PAGE_SIZE = 10;
const Users: FC = () => {
  const { t } = useTranslation('translation', { keyPrefix: 'admin.users' });

  const [urlSearchParams, setUrlSearchParams] = useSearchParams();
  const curFilter = urlSearchParams.get('filter') || UserFilterKeys[0];
  const curPage = Number(urlSearchParams.get('page') || '1');
  const curQuery = urlSearchParams.get('query') || '';
  const currentUser = loggedUserInfoStore((state) => state.user);
  const Toast = useToast();
  const {
    data,
    isLoading,
    mutate: refreshUsers,
  } = useQueryUsers({
    page: curPage,
    page_size: PAGE_SIZE,
    query: curQuery,
    ...(curFilter === 'all'
      ? {}
      : curFilter === 'staff'
      ? { staff: true }
      : { status: curFilter }),
  });
  const changeModal = useChangeModal({
    callback: refreshUsers,
  });

  const changeUserRoleModal = useChangeUserRoleModal({
    callback: refreshUsers,
  });

  const handleAction = (type, user) => {
    const { user_id, status, role_id, username } = user;
    if (username === currentUser.username) {
      Toast.onShow({
        msg: t('fobidden_operate_self', { keyPrefix: 'toast' }),
        variant: 'warning',
      });
      return;
    }
    if (type === 'status') {
      changeModal.onShow({
        id: user_id,
        type: status,
      });
    }

    if (type === 'role') {
      changeUserRoleModal.onShow({
        id: user_id,
        role_id,
      });
    }
  };

  const handleFilter = (e) => {
    urlSearchParams.set('query', e.target.value);
    urlSearchParams.delete('page');
    setUrlSearchParams(urlSearchParams);
  };
  return (
    <>
      <h3 className="mb-4">{t('title')}</h3>
      <div className="d-flex justify-content-between align-items-center mb-3">
        <QueryGroup
          data={UserFilterKeys}
          currentSort={curFilter}
          sortKey="filter"
          i18nKeyPrefix="admin.users"
        />

        <Form.Control
          size="sm"
          value={curQuery}
          onChange={handleFilter}
          placeholder={t('filter.placeholder')}
          style={{ width: '12.25rem' }}
        />
      </div>
      <Table>
        <thead>
          <tr>
            <th>{t('name')}</th>
            {/* <th style={{ width: '12%' }}>{t('reputation')}</th> */}
            <th style={{ width: '20%' }}>{t('email')}</th>
            <th className="text-nowrap" style={{ width: '15%' }}>
              {t('created_at')}
            </th>
            {(curFilter === 'deleted' || curFilter === 'suspended') && (
              <th className="text-nowrap" style={{ width: '15%' }}>
                {curFilter === 'deleted' ? t('delete_at') : t('suspend_at')}
              </th>
            )}

            <th style={{ width: '12%' }}>{t('status')}</th>
            {/* <th style={{ width: '12%' }}>{t('role')}</th> */}
            {curFilter !== 'deleted' ? (
              <th style={{ width: '8%' }} className="text-end">
                {t('action')}
              </th>
            ) : null}
          </tr>
        </thead>
        <tbody className="align-middle">
          {data?.list.map((user) => {
            return (
              <tr key={user.user_id}>
                <td>
                  <BaseUserCard
                    data={user}
                    className="fs-6"
                    avatarSize="32px"
                    avatarSearchStr="s=48"
                    avatarClass="me-2"
                    showReputation={false}
                  />
                </td>
                {/* <td>{user.rank}</td> */}
                <td className="text-break">{user.e_mail}</td>
                <td>
                  <FormatTime time={user.created_at} />
                </td>
                {curFilter === 'suspended' && (
                  <td className="text-nowrap">
                    <FormatTime time={user.suspended_at} />
                  </td>
                )}
                {curFilter === 'deleted' && (
                  <td className="text-nowrap">
                    <FormatTime time={user.deleted_at} />
                  </td>
                )}
                <td>
                  <span className={classNames('badge', bgMap[user.status])}>
                    {t(user.status)}
                  </span>
                </td>
                {/* <td> */}
                {/*  <span className="badge text-bg-light"> */}
                {/*    {t(user.role_name)} */}
                {/*  </span> */}
                {/* </td> */}
                {curFilter !== 'deleted' ? (
                  <td className="text-end">
                    <Dropdown>
                      <Dropdown.Toggle variant="link" className="no-toggle">
                        <Icon name="three-dots-vertical" />
                      </Dropdown.Toggle>
                      <Dropdown.Menu>
                        {/* <Dropdown.Item>{t('set_new_password')}</Dropdown.Item> */}
                        <Dropdown.Item
                          onClick={() => handleAction('status', user)}>
                          {t('change_status')}
                        </Dropdown.Item>
                        {/* <Dropdown.Item */}
                        {/*  onClick={() => handleAction('role', user)}> */}
                        {/*  {t('change_role')} */}
                        {/* </Dropdown.Item> */}
                        {/* <Dropdown.Divider />
                        <Dropdown.Item>{t('show_logs')}</Dropdown.Item> */}
                      </Dropdown.Menu>
                    </Dropdown>

                    {/* {user.status !== 'deleted' && (
                      <Button
                        className="p-0 btn-no-border"
                        variant="link"
                        onClick={() => handleClick(user)}>
                        {t('change')}
                      </Button>
                    )} */}
                  </td>
                ) : null}
              </tr>
            );
          })}
        </tbody>
      </Table>
      {Number(data?.count) <= 0 && !isLoading && <Empty />}
      <div className="mt-4 mb-2 d-flex justify-content-center">
        <Pagination
          currentPage={curPage}
          totalSize={data?.count || 0}
          pageSize={PAGE_SIZE}
        />
      </div>
    </>
  );
};

export default Users;
