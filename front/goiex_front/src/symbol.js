import React, {useState, useEffect} from 'react';
import 'bootstrap/dist/css/bootstrap.css';

import {Button, Card, Row, Col} from 'react-bootstrap'

const Symbol = ({symbolData}) => {

    return (
        <Card>
            <Row>
                <Col>{ symbolData !== undefined && symbolData.symbol}</Col>
                <Col>{ symbolData !== undefined && symbolData.date}</Col>
                <Col>{ symbolData !== undefined && symbolData.currency}</Col>
                <Col>{ symbolData !== undefined && symbolData.exchange}</Col>
            </Row>
        </Card>
    )

}

export default Symbol