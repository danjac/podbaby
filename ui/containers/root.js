import React, { PropTypes } from 'react';
import { Provider } from 'react-redux';
import DocumentTitle from 'react-document-title';
import { getTitle } from './utils';


export default class Root extends React.Component {

  render() {
    const { store, routes } = this.props;

    return (
      <DocumentTitle title={getTitle()}>
        <div>
        <Provider store={store}>
          {routes}
        </Provider>
        </div>
      </DocumentTitle>
    );
  }

}

Root.propTypes = {
  store: PropTypes.object.isRequired,
  routes: PropTypes.object.isRequired,
};
